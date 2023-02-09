package cdn

import (
	"context"
	"github.com/carlmjohnson/requests"
	"strings"
	"sync"
)

// CdnFetcher
//
//	CdnFetcher
//	@Description:
type CdnFetcher struct {
}

type CdnApiResponse struct {
	Ipv4Range []string `json:"ipv_4_range"`
	Ipv6Range []string `json:"ipv_6_range"`
}

type GcoreCdnApiResp struct {
	Addresses   []string `json:"addresses"`
	AddressesV6 []string `json:"addresses_v6"`
}

type CloudfrontApiResp struct {
	CloudfrontGlobalIPList       []string `json:"CLOUDFRONT_GLOBAL_IP_LIST"`
	CloudfrontRegionalEdgeIPList []string `json:"CLOUDFRONT_REGIONAL_EDGE_IP_LIST"`
}

// FetchCloudFlare
//
//	@Description: 获取cloudflare家的cdn ip range
//	@receiver receiver
//	@return CdnApiResponse
func (receiver CdnFetcher) FetchCloudFlare() CdnApiResponse {
	var apiUrl4 = "https://www.cloudflare.com/ips-v4"
	var apiUrl6 = "https://www.cloudflare.com/ips-v6"
	ctx := context.Background()
	var data4 string
	var data6 string
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := requests.
			URL(apiUrl4).
			ToString(&data4).
			Fetch(ctx)
		HandleError(err)
		// fmt.Printf("%v\n", data4)
	}()

	go func() {
		defer wg.Done()
		err := requests.
			URL(apiUrl6).
			ToString(&data6).
			Fetch(ctx)
		HandleError(err)
		// fmt.Printf("%v\n", data6)
	}()
	wg.Wait()
	splitIp4 := strings.Split(data4, "\n")
	splitIp6 := strings.Split(data6, "\n")
	return CdnApiResponse{
		Ipv4Range: splitIp4,
		Ipv6Range: splitIp6,
	}
}

// FetchGcore
//
//	@Description: 获取Gcore家的cdn range
//	@receiver receiver
//	@return CdnApiResponse
func (receiver CdnFetcher) FetchGcore() CdnApiResponse {
	var apiUrl = "https://api.gcorelabs.com/cdn/public-net-list"
	ctx := context.Background()
	var data = GcoreCdnApiResp{}
	err := requests.
		URL(apiUrl).
		ToJSON(&data).
		Fetch(ctx)
	HandleError(err)
	// fmt.Printf("%v\n", data)
	return CdnApiResponse{
		Ipv4Range: data.Addresses,
		Ipv6Range: data.AddressesV6,
	}
}

// FetchCloudfront
//
//	@Description: 获取cloudfront cdn ip range
//	@receiver receiver
//	@return CdnApiResponse
func (receiver CdnFetcher) FetchCloudfront() CdnApiResponse {
	ctx := context.Background()
	var apiUrl = "https://d7uri8nf7uskq.cloudfront.net/tools/list-cloudfront-ips"
	var data = CloudfrontApiResp{}
	err := requests.
		URL(apiUrl).
		ToJSON(&data).
		Fetch(ctx)
	HandleError(err)
	// fmt.Printf("%v\n", data)

	return CdnApiResponse{
		Ipv4Range: data.CloudfrontGlobalIPList,
		Ipv6Range: nil,
	}
}

func HandleError(err error) {
	if err != nil {
		panic("程序当前运行出错: " + err.Error())
	}
}
