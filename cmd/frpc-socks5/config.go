package main

var (
	socks5User     *string = strPtr("xxxaccesskey")
	socks5Pass     *string = strPtr("xxxaccesssecret")
	frpsToken      *string = strPtr("abcd.1234")
	frpsAddr       *string = strPtr("43.155.9.26")
	frpsPort       *int    = intPtr(7000)
	deviceId       *string = strPtr("windows071101")
	frpsRemotePort *int    = intPtr(10009)
)

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
