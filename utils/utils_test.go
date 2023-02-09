package utils

import "testing"

func TestIpv4_Ip2lon(t *testing.T) {
	ipv4 := Ipv4Tool{}
	ip := ipv4.Lon2ip(12)
	println(ip)
}

func TestIpv4_Lon2ip(t *testing.T) {
	ipv4 := Ipv4Tool{}
	lon := ipv4.Ip2lon("192.168.0.1")
	println(lon)
}

func TestIpv4_Random(t *testing.T) {
	ipv4 := Ipv4Tool{}
	random := ipv4.Random("192.168.0.1", 24)
	println(random)
}
