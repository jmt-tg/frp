package main

import (
	"context"
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config/types"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/util/log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// go run . --socks5-user=xxxaccesskey --socks5-pass=xxxaccesssecret --frps-token=abcd.1234 --frps-addr=45.194.33.6 --frps-port=7000 --device-id=macos --frps-remote-port=11080

//var socks5User = flag.String("socks5-user", "xxxaccesskey", "socks5 user")
//var socks5Pass = flag.String("socks5-pass", "xxxaccesssecret", "socks5 pass")
//var frpsToken = flag.String("frps-token", "", "frps token")
//var frpsAddr = flag.String("frps-addr", "", "frps addr")
//var frpsPort = flag.Int("frps-port", 7000, "frps port")
//var deviceId = flag.String("device-id", "windows", "device id")
//var frpsRemotePort = flag.Int("frps-remote-port", 11080, "frps remote port")

func main() {
	//flag.Parse()
	socks5Port := freePort()
	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{socks5.UserPassAuthenticator{
			Credentials: socks5.StaticCredentials{
				*socks5User: *socks5Pass,
			},
		}},
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := server.ListenAndServe(
			"tcp",
			fmt.Sprintf("127.0.0.1:%d", socks5Port),
		); err != nil {
			panic(err)
		}
	}()

	cfg := &v1.ClientCommonConfig{
		Auth: v1.AuthClientConfig{
			Method: "token",
			Token:  *frpsToken,
		},
		ServerAddr: *frpsAddr,
		ServerPort: *frpsPort,
		Log: v1.LogConfig{
			To:                "console",
			Level:             "error",
			MaxDays:           0,
			DisablePrintColor: false,
		},
		Transport: v1.ClientTransportConfig{
			Protocol: "tcp",
		},
	}
	proxyCfgs := []v1.ProxyConfigurer{
		&v1.TCPProxyConfig{
			ProxyBaseConfig: v1.ProxyBaseConfig{
				Name: "local-socks5-" + *deviceId,
				Type: "tcp",
				ProxyBackend: v1.ProxyBackend{
					LocalIP:   "127.0.0.1",
					LocalPort: socks5Port,
					Plugin:    v1.TypedClientPluginOptions{},
				},
				Transport: v1.ProxyTransport{
					UseEncryption:        false,
					UseCompression:       false,
					BandwidthLimit:       types.BandwidthQuantity{},
					BandwidthLimitMode:   "client",
					ProxyProtocolVersion: "",
				},
			},
			RemotePort: *frpsRemotePort,
		},
	}
	var visitorCfgs []v1.VisitorConfigurer
	warning, err := validation.ValidateAllClientConfig(cfg, proxyCfgs, visitorCfgs)
	if warning != nil {
		fmt.Printf("WARNING: %v\n", warning)
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("frpc start with socks5 server at: %d\n", socks5Port)
	err = startFrpc(cfg, proxyCfgs, visitorCfgs)
	if err != nil {
		panic(err)
	}
}

func startFrpc(
	cfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer,
) error {
	log.InitLogger(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor)
	svr, err := client.NewService(client.ServiceOptions{
		Common:      cfg,
		ProxyCfgs:   proxyCfgs,
		VisitorCfgs: visitorCfgs,
	})
	if err != nil {
		return err
	}

	shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
	// Capture the exit signal if we use kcp or quic.
	if shouldGracefulClose {
		go handleTermSignal(svr)
	}
	return svr.Run(context.Background())
}

func handleTermSignal(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}

func freePort() int {
	// 获取空闲端口
	var a *net.TCPAddr
	var err error
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port
		}
	}
	panic(err)
}

// export https_proxy=socks5://xxxaccesskey:xxxaccesssecret@127.0.0.1:1080 http_proxy=socks5://xxxaccesskey:xxxaccesssecret@127.0.0.1:1080 all_proxy=socks5://xxxaccesskey:xxxaccesssecret@127.0.0.1:1080
// curl https://www.google.com

// export https_proxy=socks5://xxxaccesskey:xxxaccesssecret@45.194.33.6:11080 http_proxy=socks5://xxxaccesskey:xxxaccesssecret@45.194.33.6:11080 all_proxy=socks5://xxxaccesskey:xxxaccesssecret@45.194.33.6:11080
// curl https://www.google.com
