package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
)

const IpGeoLocationApiKey = "71a993e55ea64df29c3caa7c094f7099"

// https://ipgeolocation.io/ data
// 1k request for per day

type IpgeolocationInfo struct {
	IP             string   `json:"ip"`
	Hostname       string   `json:"hostname"`
	ContinentCode  string   `json:"continent_code"`
	ContinentName  string   `json:"continent_name"`
	CountryCode2   string   `json:"country_code2"`
	CountryCode3   string   `json:"country_code3"`
	CountryName    string   `json:"country_name"`
	CountryCapital string   `json:"country_capital"`
	StateProv      string   `json:"state_prov"`
	District       string   `json:"district"`
	City           string   `json:"city"`
	Zipcode        string   `json:"zipcode"`
	Latitude       string   `json:"latitude"`
	Longitude      string   `json:"longitude"`
	IsEu           bool     `json:"is_eu"`
	CallingCode    string   `json:"calling_code"`
	CountryTld     string   `json:"country_tld"`
	Languages      string   `json:"languages"`
	CountryFlag    string   `json:"country_flag"`
	GeonameID      string   `json:"geoname_id"`
	Isp            string   `json:"isp"`
	ConnectionType string   `json:"connection_type"`
	Organization   string   `json:"organization"`
	Asn            string   `json:"asn"`
	Currency       Currency `json:"currency"`
	TimeZone       TimeZone `json:"time_zone"`
}
type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
type TimeZone struct {
	Name            string  `json:"name"`
	Offset          int     `json:"offset"`
	CurrentTime     string  `json:"current_time"`
	CurrentTimeUnix float64 `json:"current_time_unix"`
	IsDst           bool    `json:"is_dst"`
	DstSavings      int     `json:"dst_savings"`
}

// GetIpgeolocationInfo
//
//	@Description: 使用ipgeolocation 获取geoip信息
//	@param ipString
//	@return *IpgeolocationInfo
//	@return error
func GetIpgeolocationInfo(ipString string) (*IpgeolocationInfo, error) {
	normalIpv4Address := CheckStrIsIpAddress(ipString)
	if !normalIpv4Address {
		return nil, errors.New("args not a valid ipStr address: " + ipString)
	}
	var requestUrl = fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", IpGeoLocationApiKey, ipString)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Println("get location by ipgeolocation error: ", err)
		return nil, errors.New("get location by ipgeolocation error: " + err.Error())
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "zh,en;q=0.9,zh-TW;q=0.8,zh-CN;q=0.7,ja;q=0.6")
	req.Header.Add("Cache-Control", "no-cache")
	//req.Header.Add("Origin", "http://ip-api.com")
	req.Header.Add("Pragma", "no-cache")
	//req.Header.Add("Referer", "http://ip-api.com/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	var ipGeoInfo IpgeolocationInfo
	err = json.NewDecoder(resp.Body).Decode(&ipGeoInfo)
	if err != nil {
		log.Println("get location by ipgeolocation error: ", err)
		return nil, errors.New("get location by ipgeolocation error: " + err.Error())
	}
	return &ipGeoInfo, nil
}

// GetLocInfoShort
//
//	@Description: 获取地理位置信息
//	@receiver receiver
//	@return string
func (receiver *IpgeolocationInfo) GetLocInfoShort() string {
	if receiver == nil {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", receiver.CountryCode3, receiver.City, receiver.Organization)
}

// CheckStrIsIpAddress
//
//	@Description: 判断str是否为合格的ip str
//	@param str
//	@return bool
func CheckStrIsIpAddress(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil
}
