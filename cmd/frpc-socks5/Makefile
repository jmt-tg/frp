exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ~/sdk/go-latest/bin/go build -ldflags "-s -w -H=windowsgui" -o frpc-socks5.exe
	echo "编译完成，大小：" && ls -lh frpc-socks5.exe | awk '{print $$5}'

test-windows:
	.\frpc-socks5.exe --socks5-user=xxxaccesskey --socks5-pass=xxxaccesssecret --frps-token=abcd.1234 --frps-addr=45.194.33.6 --frps-port=7000 --device-id=macos --frps-remote-port=11080

mac-bin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 ~/sdk/go-latest/bin/go build -ldflags "-s -w" -o frpc-socks5
	echo "编译完成，大小：" && ls -lh frpc-socks5 | awk '{print $$5}'

test-macos:
	chmod +x ./frpc-socks5
	#./frpc-socks5 --socks5-user=xxxaccesskey --socks5-pass=xxxaccesssecret --frps-token=abcd.1234 --frps-addr=45.194.33.6 --frps-port=7000 --device-id=macos --frps-remote-port=11080
	./frpc-socks5
	#export https_proxy=socks5://xxxaccesskey:xxxaccesssecret@45.194.33.6:11080 http_proxy=socks5://xxxaccesskey:xxxaccesssecret@45.194.33.6:11080 all_proxy=socks5://xxxaccesskey:xxxaccesssecret@45.194.33.6:11080
    #curl https://www.google.com
    #curl ifconfig.me
