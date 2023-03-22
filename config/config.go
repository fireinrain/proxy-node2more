package config

// AllConfig 配置枚举
type AllConfig struct {
	//输入的节点切片
	InputNodeStr []string `json:"input_node_str"`
	//cdn提供商
	CDNName CdnProvider `json:"cdn_name"`
	//自定义cdn ip
	CustomCDNIp []string `json:"custom_cdn_ip"`
	//获取方式
	GetMethodName GetMethod `json:"get_method_name"`
	//获取的节点数
	WantedNodeNum int `json:"wanted_node_num"`
	//输出的节点切片
	OutPutNodeList []string `json:"out_put_node_list"`
}

// CdnProvider cdn提供商枚举
type CdnProvider int

const (
	CDNCloudflare CdnProvider = iota
	CDNCloudFront
	CDNGcore
	CDNOther
)

// GetMethod 获取方式
type GetMethod int

const (
	GetMethodSequance GetMethod = iota
	GetMethodRandom
)

// VmessInfo
//
//	VmessInfo
//	@Description: vmess 信息
type VmessInfo struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  string `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	Sni  string `json:"sni"`
	Alpn string `json:"alpn"`
}

func (receiver VmessInfo) CloneNew() VmessInfo {
	return VmessInfo{
		V:    receiver.V,
		Ps:   receiver.Ps,
		Add:  receiver.Add,
		Port: receiver.Port,
		ID:   receiver.ID,
		Aid:  receiver.Aid,
		Scy:  receiver.Scy,
		Net:  receiver.Net,
		Type: receiver.Type,
		Host: receiver.Host,
		Path: receiver.Path,
		TLS:  receiver.TLS,
		Sni:  receiver.Sni,
		Alpn: receiver.Alpn,
	}

}
