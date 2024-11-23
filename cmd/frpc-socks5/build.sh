#!/bin/bash


function build() {
    frpsRemotePort=$1
    echo "开始编译frpc-socks5.exe, frpsRemotePort=${frpsRemotePort}"
    rm -f config.go
    cp config.go.tmpl config.go
    sed -i "" "s#FRPS_REMOTE_PORT#${frpsRemotePort}#g" config.go
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ~/sdk/go-latest/bin/go build -ldflags "-s -w -H=windowsgui" -o "frpc-socks5-${frpsRemotePort}.exe"
    echo "编译完成，大小：" && ls -lh "frpc-socks5-${frpsRemotePort}.exe" | awk '{print $$5}'
}


build 10001
build 10002
build 10003
build 10004
build 10005
build 10006
build 10007
build 10008
build 10009

zip frpc-socks5.zip frpc-socks5-*.exe
