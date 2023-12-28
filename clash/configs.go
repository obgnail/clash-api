package clash

import (
	"fmt"
	"github.com/juju/errors"
)

type Tun struct {
	Enable              bool     `json:"enable"`
	Device              string   `json:"device"`
	Stack               string   `json:"stack"`
	DNSHijack           []string `json:"dns-hijack"`
	AutoRoute           bool     `json:"auto-route"`
	AutoDetectInterface bool     `json:"auto-detect-interface"`
	MTU                 int      `json:"mtu"`
	GSOMaxSize          int      `json:"gso-max-size"`
	Inet4Address        []string `json:"inet4-address"`
	StrictRoute         bool     `json:"strict-route"`
	FileDescriptor      int      `json:"file-descriptor"`
}

type MuxOption struct {
	Brutal struct {
		Enabled bool `json:"enabled"`
	} `json:"brutal"`
}

type TuicServer struct {
	Enable      bool      `json:"enable"`
	Listen      string    `json:"listen"`
	Certificate string    `json:"certificate"`
	PrivateKey  string    `json:"private-key"`
	MuxOption   MuxOption `json:"mux-option"`
}

type Configs struct {
	Port                    int               `json:"port"`
	SocksPort               int               `json:"socks-port"`
	RedirPort               int               `json:"redir-port"`
	TProxyPort              int               `json:"tproxy-port"`
	MixedPort               int               `json:"mixed-port"`
	Tun                     Tun               `json:"tun"`
	TuicServer              TuicServer        `json:"tuic-server"`
	SSConfig                string            `json:"ss-config"`
	VmessConfig             string            `json:"vmess-config"`
	Authentication          interface{}       `json:"authentication"`
	SkipAuthPrefixes        []string          `json:"skip-auth-prefixes"`
	LanAllowedIPs           []string          `json:"lan-allowed-ips"`
	LanDisallowedIPs        interface{}       `json:"lan-disallowed-ips"`
	AllowLAN                bool              `json:"allow-lan"`
	BindAddress             string            `json:"bind-address"`
	InboundTFO              bool              `json:"inbound-tfo"`
	InboundMPTCP            bool              `json:"inbound-mptcp"`
	Mode                    string            `json:"mode"`
	UnifiedDelay            bool              `json:"UnifiedDelay"`
	LogLevel                string            `json:"log-level"`
	IPv6                    bool              `json:"ipv6"`
	InterfaceName           string            `json:"interface-name"`
	GeoxURL                 map[string]string `json:"geox-url"`
	GeoAutoUpdate           bool              `json:"geo-auto-update"`
	GeoUpdateInterval       int               `json:"geo-update-interval"`
	GeodataMode             bool              `json:"geodata-mode"`
	GeodataLoader           string            `json:"geodata-loader"`
	TCPConcurrent           bool              `json:"tcp-concurrent"`
	FindProcessMode         string            `json:"find-process-mode"`
	Sniffing                bool              `json:"sniffing"`
	GlobalClientFingerprint string            `json:"global-client-fingerprint"`
	GlobalUA                string            `json:"global-ua"`
}

func GetConfigs() (*Configs, error) {
	configs := &Configs{}
	err := UnmarshalRequest("get", "/configs", nil, nil, &configs)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return configs, nil
}

// EnableConfigs 这个接口不会影响 external-controller 和 secret 的值
func EnableConfigs(path string) error {
	headers := map[string]string{"Content-Type": "application/json"}
	body := map[string]interface{}{"path": path}

	code, content, err := EasyRequest("put", "/configs?force=true", headers, body)
	if err != nil {
		return errors.Trace(err)
	}
	if code != 200 {
		return fmt.Errorf("unknown error: %s", string(content))
	}
	return nil
}

// SetConfigs 更新基本配置，传入需要修改的配置即可，传入的数据需以 json 格式传入
//
// 命令行调用示例：curl ${controller-api}/configs -X PATCH -d '{"mixed-port": 7890}'
func SetConfigs(body map[string]interface{}) error {
	headers := map[string]string{"Content-Type": "application/json"}

	code, content, err := EasyRequest("patch", "/configs", headers, body)
	if err != nil {
		return errors.Trace(err)
	}
	if code != 200 {
		return fmt.Errorf("unknown error: %s", string(content))
	}
	return nil
}

// EnableGeo 更新 GEO 数据库,因更新后会自动重载一次配置
func EnableGeo() error {
	code, content, err := EasyRequest("post", "/configs/geo", nil, nil)
	if err != nil {
		return errors.Trace(err)
	}
	if code != 200 {
		return fmt.Errorf("unknown error: %s", string(content))
	}
	return nil
}
