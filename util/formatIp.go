package util

import "net"

func Ip2long(s string) (ret int64) {
	bip := ([]byte)(net.ParseIP(s).To4())
	return (int64)(bip[0])*(1<<24) + (int64)(bip[1])*(1<<16) + (int64)(bip[2])*(1<<8) + (int64)(bip[3])
}
