package utils

import (
	"fmt"
	"math"
	"math/rand"
	"net/netip"
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
