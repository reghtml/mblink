package GoMiniblink

// ProxyType 代理服务器类型
type ProxyType int

const (
	ProxyType_NONE ProxyType = iota
	ProxyType_HTTP
	ProxyType_SOCKS4
	ProxyType_SOCKS4A
	ProxyType_SOCKS5
	ProxyType_SOCKS5HOSTNAME
)

// ProxyInfo 代理服务器配置信息
type ProxyInfo struct {
	Type     ProxyType
	HostName string
	Port     int
	UserName string
	Password string
}
