package main

var (
	socks5User     *string = strPtr("xxxaccesskey")
	socks5Pass     *string = strPtr("xxxaccesssecret")
	frpsToken      *string = strPtr("abcd.1234")
	frpsAddr       *string = strPtr("45.194.33.6")
	frpsPort       *int    = intPtr(7000)
	deviceId       *string = strPtr("windows")
	frpsRemotePort *int    = intPtr(11000)
)

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
