package utils

import (
	"proxy-node2more/config"
	"testing"
)

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

func TestCaculateNodesResult(t *testing.T) {
	config := &config.AllConfig{
		InputNodeStr:   "vmess://ew0KICAidiI6ICIyIiwNCiAgInBzIjogImRkcGMtdm1lc3MtY2RuLWNmdCIsDQogICJhZGQiOiAiZDFrOGJld2p5bmt0em0uY2xvdWRmcm9udC5uZXQiLA0KICAicG9ydCI6ICI0NDMiLA0KICAiaWQiOiAiZjc3Zjk3M2MtNDRlYi00ZDFlLWZiZjAtMGZhYzY2ZGZjYzQzIiwNCiAgImFpZCI6ICIwIiwNCiAgInNjeSI6ICJhdXRvIiwNCiAgIm5ldCI6ICJ3cyIsDQogICJ0eXBlIjogIm5vbmUiLA0KICAiaG9zdCI6ICIiLA0KICAicGF0aCI6ICIvZmlyZSIsDQogICJ0bHMiOiAidGxzIiwNCiAgInNuaSI6ICIiLA0KICAiYWxwbiI6ICIiDQp9",
		CDNName:        0,
		CustomCDNIp:    nil,
		GetMethodName:  0,
		WantedNodeNum:  10,
		OutPutNodeList: nil,
	}
	result, err := CaculateNodesResult(config)
	HandleError(err)
	println(result)
}
