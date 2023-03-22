package utils

import (
	"fmt"
	"proxy-node2more/config"
	"regexp"
	"strings"
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
		InputNodeStr:   []string{"vmess://ew0KICAidiI6ICIyIiwNCiAgInBzIjogImRkcGMtdm1lc3MtY2RuLWNmdCIsDQogICJhZGQiOiAiZDFrOGJld2p5bmt0em0uY2xvdWRmcm9udC5uZXQiLA0KICAicG9ydCI6ICI0NDMiLA0KICAiaWQiOiAiZjc3Zjk3M2MtNDRlYi00ZDFlLWZiZjAtMGZhYzY2ZGZjYzQzIiwNCiAgImFpZCI6ICIwIiwNCiAgInNjeSI6ICJhdXRvIiwNCiAgIm5ldCI6ICJ3cyIsDQogICJ0eXBlIjogIm5vbmUiLA0KICAiaG9zdCI6ICIiLA0KICAicGF0aCI6ICIvZmlyZSIsDQogICJ0bHMiOiAidGxzIiwNCiAgInNuaSI6ICIiLA0KICAiYWxwbiI6ICIiDQp9"},
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

func TestTrojanLink(t *testing.T) {
	sampleNode := "trojan://4bd8ab61-7e87-4ee6-be58-fe14fc62e6c0@ca1.trojanvh.xyz:80?security=tls&sni=ca1.trojanvh.xyz&type=tcp&headerType=none#org-org.org_Relay_-%F0%9F%87%A8%F0%9F%87%A6CA_36"
	re := regexp.MustCompile(`(@)(.*?)(:)(.*?)(\?)`)
	subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
	fmt.Println(subStrPart1)

}

func TestTrojanLink2(t *testing.T) {
	sampleNode := "trojan://BQXJHGqe@trojan.wefuckgfw.tk:443#trojan.wefuckgfw.tk%3A443"
	var re *regexp.Regexp
	if strings.Contains(sampleNode, "?") {
		re := regexp.MustCompile(`(@)(.*?)(:)(.*?)(\?)`)
		subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
		fmt.Println(subStrPart1)
	} else {
		re := regexp.MustCompile(`(@)(.*?)(:)(\d+)(.*?)`)
		subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
		subStrPart2 := re.FindStringSubmatch(sampleNode)[2]
		subStrPart3 := re.FindStringSubmatch(sampleNode)[3]
		subStrPart4 := re.FindStringSubmatch(sampleNode)[4]
		subStrPart5 := re.FindStringSubmatch(sampleNode)[5]

		fmt.Println(subStrPart1)
		fmt.Println(subStrPart2)
		fmt.Println(subStrPart3)
		fmt.Println(subStrPart4)
		fmt.Println(subStrPart5)

	}
	println(re)

}
