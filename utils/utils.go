package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/netip"
	"proxy-node2more/cdn"
	"proxy-node2more/config"
	"proxy-node2more/location"
	"regexp"
	"strconv"
	"strings"
)

type Ipv4Tool struct {
}

// Random
//
//	@Description: 随机获取CIDR 中的一个ip
//
// 192.168.0.1/24 ----> 192.168.0.0 - 192.168.0.255
//
//	@receiver receiver
//	@param subset
//	@param mask
//	@return string
func (receiver *Ipv4Tool) Random(subset string, mask int) string {
	var randomIp = math.Floor(rand.Float64()*math.Pow(2, 32-float64(mask)) + 1)
	i := receiver.Ip2lon(subset) | int(randomIp)
	return receiver.Lon2ip(i)
}

// Ip2lon
//
//	@Description: ip地址转化为一个int整数
//	@receiver receiver
//	@param address
//	@return int
func (receiver *Ipv4Tool) Ip2lon(address string) int {
	var result int = 0
	splits := strings.Split(address, ".")
	for _, value := range splits {
		result <<= 8
		atoi, _ := strconv.Atoi(value)
		result += atoi
	}
	return int(uint32(result) >> 0)

}

// Lon2ip
//
//	@Description: 一个大整数转化为ip地址
//	@receiver receiver
//	@param lon
//	@return string
func (receiver *Ipv4Tool) Lon2ip(lon int) string {
	var a = uint32(lon) >> 24
	var b = lon >> 16 & 255
	var c = lon >> 8 & 255
	var d = lon & 255
	aa := strconv.Itoa(int(a))
	bb := strconv.Itoa(b)
	cc := strconv.Itoa(c)
	dd := strconv.Itoa(d)
	ipStrs := []string{}
	ipStrs = append(ipStrs, aa)
	ipStrs = append(ipStrs, bb)
	ipStrs = append(ipStrs, cc)
	ipStrs = append(ipStrs, dd)

	return strings.Join(ipStrs, ".")
}

// CIDR2IPS
//
//	@Description: 从cidr计算所有包含的ip
//	@param cidr
//	@return []netip.Addr
//	@return error
func CIDR2IPS(cidr string) ([]netip.Addr, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		panic(err)
	}

	var ips []netip.Addr
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr)
	}

	if len(ips) < 2 {
		return ips, nil
	}

	return ips[1 : len(ips)-1], nil
}

func HandleError(err error) {
	if err != nil {
		fmt.Println("处理出现错误: ", err.Error())
		panic(err)
	}
}

// CaculateNodesResult
//
//	@Description: 计算并替换节点
//	@param configSet
//	@return *config.AllConfig
//	@return error
func CaculateNodesResult(configSet *config.AllConfig) (*config.AllConfig, error) {
	var inputNodes = configSet.InputNodeStr
	var nodeNum = configSet.WantedNodeNum
	var ipResult = []string{}
	var giplist = []string{}
	ipv4Tool := Ipv4Tool{}
	//获取cloudflare
	cdnFetcher := cdn.GlobalCdnFetcher

	if configSet.CDNName == 0 {
		cloudFlare := cdnFetcher.FetchCloudFlare()
		giplist = cloudFlare.Ipv4Range
	}
	//Gcore
	if configSet.CDNName == 1 {
		gcore := cdnFetcher.FetchGcore()
		giplist = gcore.Ipv4Range
	}
	if configSet.CDNName == 2 {
		cloudfront := cdnFetcher.FetchCloudfront()
		giplist = cloudfront.Ipv4Range
	}
	if configSet.CDNName == 3 {
		giplist = configSet.CustomCDNIp
	}

	//生成ip列表
	//顺序
	if configSet.GetMethodName == 0 {
		//如果是自定义ip
		if configSet.CDNName == 3 {
			ipResult = giplist
		} else {
			for i, index := 0, 0; i < nodeNum; i, index = i+1, index+1 {
				if len(giplist) == index {
					index = 0
				}
				ipMaskResult := strings.Split(giplist[index], "/")
				mask, err := strconv.Atoi(ipMaskResult[1])
				HandleError(err)
				//判断生成的ip是否已经存在(去重处理)
				var randomIp string
				for {
				start:
					randomIp = ipv4Tool.Random(ipMaskResult[0], mask)
					for _, ip := range ipResult {
						if ip == randomIp {
							goto start
						}
					}
					break
				}
				ipResult = append(ipResult, randomIp)
			}

		}

	} else {
		if configSet.CDNName == 3 {
			ipResult = giplist
		} else {
			//随机
			for i := 0; i < nodeNum; i++ {
				index := rand.Intn(len(giplist))
				ipMaskResult := strings.Split(giplist[index], "/")
				mask, err := strconv.Atoi(ipMaskResult[1])
				HandleError(err)
				//判断生成的ip是否已经存在(去重处理)
				var randomIp string
				for {
				start2:
					randomIp = ipv4Tool.Random(ipMaskResult[0], mask)
					for _, ip := range ipResult {
						if ip == randomIp {
							goto start2
						}
					}
					break
				}
				ipResult = append(ipResult, randomIp)
			}
		}

	}
	//
	var ipResultGeo = make(map[string]string)
	for _, ip := range ipResult {
		//获取位置信息
		info, err := location.GetIpgeolocationInfo(ip)
		if err != nil {
			fmt.Println("Get ip geolocation info error: ", err)
		}
		infoShort := info.GetLocInfoShort()
		ipResultGeo[ip] = infoShort
	}

	var results []string
	for _, node := range inputNodes {
		replaceCDNIpSet, err := DoReplaceCDNIpSet(node, ipResultGeo)
		if err != nil {
			return &config.AllConfig{
				InputNodeStr:   nil,
				CDNName:        0,
				GetMethodName:  0,
				WantedNodeNum:  0,
				OutPutNodeList: nil,
			}, err
		}
		for _, n := range replaceCDNIpSet {
			results = append(results, n)
		}
	}
	configSet.OutPutNodeList = results
	return configSet, nil
}

func DoReplaceCDNIpSet(inputNode string, ipResultGeo map[string]string) ([]string, error) {
	ipResult := make([]string, 0, len(ipResultGeo))
	for k := range ipResultGeo {
		ipResult = append(ipResult, k)
	}
	output := inputNode

	output = strings.ReplaceAll(output, " ", "")
	output = strings.ReplaceAll(output, "\t", "")
	output = strings.ReplaceAll(output, "\n", "")

	sampleNode := output
	vmessPre := "vmess://"
	vlessPre := "vless://"
	trojanPre := "trojan://"

	if !strings.HasPrefix(sampleNode, vmessPre) && !strings.HasPrefix(sampleNode, vlessPre) && !strings.HasPrefix(sampleNode, trojanPre) {
		return nil, errors.New("仅支持vmess、vless和trojan的节点分享链接")
	}

	//对ip列表去重
	//reduceMap := make(map[string]int)
	//newSlice := []string{}
	//
	//for _, ipnode := range ipResult {
	//	if _, value := reduceMap[ipnode]; !value {
	//		reduceMap[ipnode] = 0
	//		newSlice = append(newSlice, ipnode)
	//	}
	//}
	//ipResult = newSlice
	//将cdn ip替换到输入的节点

	var nodes = []string{}
	//vmess节点处理
	if strings.HasPrefix(sampleNode, vmessPre) {
		vmessStr, err := atob(strings.Replace(sampleNode, vmessPre, "", 1))
		HandleError(err)
		//序列化
		vmessInfo := config.VmessInfo{}
		bytes := []byte(vmessStr)
		err = json.Unmarshal(bytes, &vmessInfo)
		HandleError(err)

		add := vmessInfo.Add
		vmessInfo.Host = add
		vmessInfo.Aid = 0
		for i := 0; i < len(ipResult); i++ {
			var newNode = vmessInfo.CloneNew()
			locationInfo := ipResultGeo[ipResult[i]]
			newNode.Ps = locationInfo + "@" + newNode.Ps
			newNode.Add = ipResult[i]
			replacedNode, err := json.Marshal(newNode)
			HandleError(err)
			s := string(replacedNode)
			s2 := btoa(s)
			nodes = append(nodes, vmessPre+s2+"\r")
		}
		return nodes, nil
	}
	//vless协议
	if strings.HasPrefix(sampleNode, vlessPre) {
		//vless:  vless://9bc0eacc-68f3-4562-15bedad6f6ef@a.b.c:539?type=tcp&security=tls&sni=b.a.tk&flow=xtls-rprx-direct#abc-vless-1
		//vless://85b3a84e-a939-4be0-b35c-80e8c9b0f1c8@xray.wefuckgfw.tk:2083?encryption=none&security=tls&type=ws&host=&path=/4TLoIulw/#xray.wefuckgfw.tk%3A2083
		re := regexp.MustCompile(`@(.*?):`)

		nodeHost := re.FindStringSubmatch(sampleNode)[1]

		if strings.Index(sampleNode, "host=") != -1 {
			re = regexp.MustCompile(`(host=)(.*?)(&)`)

			subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
			subStrPart3 := re.FindStringSubmatch(sampleNode)[3]
			sampleNode = re.ReplaceAllString(sampleNode, subStrPart1+nodeHost+subStrPart3)
		} else {
			var subStrPart5 string
			if strings.Contains(sampleNode, "?") {
				re = regexp.MustCompile(`(@)(.*?)(:)(.*?)(\?)`)
				subStrPart5 = re.FindStringSubmatch(sampleNode)[5]

			} else {
				re = regexp.MustCompile(`(@)(.*?)(:)(\d+)(.*?)`)
				subStrPart5 = re.FindStringSubmatch(sampleNode)[5] + "?"
			}
			subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
			subStrPart2 := re.FindStringSubmatch(sampleNode)[2]

			subStrPart3 := re.FindStringSubmatch(sampleNode)[3]

			subStrPart4 := re.FindStringSubmatch(sampleNode)[4]

			//subStrPart5 := re.FindStringSubmatch(sampleNode)[5]
			sampleNode = re.ReplaceAllString(sampleNode, subStrPart1+subStrPart2+subStrPart3+subStrPart4+subStrPart5+"host="+nodeHost+"&")

		}

		for _, ip := range ipResult {
			re = regexp.MustCompile(`(@)(.*?)(:)`)
			subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
			subStrPart3 := re.FindStringSubmatch(sampleNode)[3]
			nodeNameIndex := strings.LastIndex(sampleNode, "#")
			nodeName := sampleNode[nodeNameIndex:]
			locationInfo := ipResultGeo[ip]
			sampleNode = strings.ReplaceAll(sampleNode, nodeName, locationInfo+"@"+nodeName)
			nodes = append(nodes, re.ReplaceAllString(sampleNode, subStrPart1+ip+subStrPart3)+"\r")
		}
		return nodes, nil
	}
	//trojan协议
	if strings.HasPrefix(sampleNode, trojanPre) {
		//trojan  trojan://aNbwlRsdsasdasr8N@a.b.tk:48857?type=tcp&security=tls&sni=a.b.tk&flow=xtls-rprx-direct#a.b.tk-trojan-2
		re := regexp.MustCompile(`@(.*?):`)

		nodeHost := re.FindStringSubmatch(sampleNode)[1]

		if strings.Index(sampleNode, "host=") != -1 {
			re = regexp.MustCompile(`(host=)(.*?)(&)`)

			subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
			subStrPart3 := re.FindStringSubmatch(sampleNode)[3]
			sampleNode = re.ReplaceAllString(sampleNode, subStrPart1+nodeHost+subStrPart3)
		} else {
			//是否有额外参数
			var subStrPart5 string
			if strings.Contains(sampleNode, "?") {
				re = regexp.MustCompile(`(@)(.*?)(:)(.*?)(\?)`)
				subStrPart5 = re.FindStringSubmatch(sampleNode)[5]

			} else {
				re = regexp.MustCompile(`(@)(.*?)(:)(\d+)(.*?)`)
				subStrPart5 = re.FindStringSubmatch(sampleNode)[5] + "?"

			}
			subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
			subStrPart2 := re.FindStringSubmatch(sampleNode)[2]

			subStrPart3 := re.FindStringSubmatch(sampleNode)[3]

			subStrPart4 := re.FindStringSubmatch(sampleNode)[4]
			//subStrPart5 := re.FindStringSubmatch(sampleNode)[5]

			sampleNode = re.ReplaceAllString(sampleNode, subStrPart1+subStrPart2+subStrPart3+subStrPart4+subStrPart5+"host="+nodeHost+"&")
		}

		for _, ip := range ipResult {
			re = regexp.MustCompile(`(@)(.*?)(:)`)
			subStrPart1 := re.FindStringSubmatch(sampleNode)[1]
			subStrPart3 := re.FindStringSubmatch(sampleNode)[3]
			nodeNameIndex := strings.LastIndex(sampleNode, "#")
			nodeName := sampleNode[nodeNameIndex:]
			locationInfo := ipResultGeo[ip]
			sampleNode = strings.ReplaceAll(sampleNode, nodeName, locationInfo+"@"+nodeName)
			nodes = append(nodes, re.ReplaceAllString(sampleNode, subStrPart1+ip+subStrPart3)+"\r")
		}
		return nodes, nil
	}

	return nil, errors.New("仅支持vmess、vless和trojan的节点分享链接")
}

// atob
//
//	@Description: base64字符串解码
//	@param s
//	@return string
//	@return error
func atob(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// btoa
//
//	@Description: 字符串编码为base64
//	@param s
//	@return string
func btoa(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
