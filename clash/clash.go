package clash

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`\[(.+?)\](.+?)lAddr=(.+?)rAddr=(.+?)mode=(.+?)rule=(.+?)proxy=(.+)`)

type LogMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (msg *LogMessage) ToLog() *Log {
	match := reg.FindAllStringSubmatch(msg.Payload, -1)
	if len(match) != 0 && len(match[0]) != 0 {
		l := &Log{
			Type:    msg.Type,
			Socket:  strings.TrimSpace(match[0][1]),
			Message: strings.TrimSpace(match[0][2]),
			LAddr:   strings.TrimSpace(match[0][3]),
			RAddr:   strings.TrimSpace(match[0][4]),
			Mode:    strings.TrimSpace(match[0][5]),
			Rule:    strings.TrimSpace(match[0][6]),
			Proxy:   strings.TrimSpace(match[0][7]),
		}
		return l
	}
	return &Log{Type: msg.Type, Error: msg.Payload}
}

type Log struct {
	Type    string `json:"type"`
	Socket  string `json:"socket"`
	Message string `json:"message"`
	LAddr   string `json:"local_addr"`
	RAddr   string `json:"remote_addr"`
	Mode    string `json:"mode"`
	Rule    string `json:"rule"`
	Proxy   string `json:"proxy"`
	Error   string `json:"error"`
}

type Traffic struct {
	Up   uint64 `json:"up"`
	Down uint64 `json:"down"`
}

type LogLevel string

const (
	LevelError   LogLevel = "error"
	LevelInfo    LogLevel = "info"
	LevelWarning LogLevel = "warning"
	LevelDebug   LogLevel = "debug"
)

func GetLogs(level LogLevel) (chan *Log, error) {
	logChan := make(chan *Log, 1024)

	headers := map[string]string{"level": string(level)}
	resp, err := Request("get", "/logs", headers, nil)
	if err != nil {
		return logChan, errors.Trace(err)
	}

	HandleStreamResp(resp, func(line []byte) (stop bool) {
		msg := &LogMessage{}
		if err := json.Unmarshal(line, msg); err != nil {
			return true
		}
		logChan <- msg.ToLog()
		return false
	})

	return logChan, nil
}

func GetTraffic(handler func(traffic *Traffic) (stop bool)) error {
	resp, err := Request("get", "/traffic", nil, nil)
	if err != nil {
		return errors.Trace(err)
	}

	HandleStreamResp(resp, func(line []byte) (stop bool) {
		traffic := &Traffic{}
		if err := json.Unmarshal(line, traffic); err != nil {
			return true
		}
		if _stop := handler(traffic); _stop {
			return true
		}
		return false
	})
	return nil
}

type Proxies struct {
	All     []string   `json:"all"`
	History []*History `json:"history"`
	Name    string     `json:"name"`
	Now     string     `json:"now"`
	Type    string     `json:"type"`
	UDP     bool       `json:"udp"`
}

type History struct {
	Time      string `json:"time"`
	Delay     int    `json:"delay"`
	MeanDelay int    `json:"meanDelay"`
}

type Proxy struct {
	History []*History `json:"history"`
	Name    string     `json:"name"`
	Type    string     `json:"type"`
	UDP     bool       `json:"udp"`
}

func GetProxies() (map[string]*Proxies, error) {
	container := struct {
		Proxies map[string]*Proxies `json:"proxies"`
	}{}
	err := UnmarshalRequest("get", "/proxies", nil, nil, &container)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return container.Proxies, nil
}

func GetProxyMessage(name string) (*Proxy, error) {
	proxy := &Proxy{}
	route := "/proxies/" + name
	err := UnmarshalRequest("get", route, nil, nil, &proxy)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return proxy, nil
}

type ProxyDelay struct {
	Delay     int    `json:"delay"`
	MeanDelay int    `json:"meanDelay"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}

func GetProxyDelay(name string, url string, timeout int) (*ProxyDelay, error) {
	proxyDelay := &ProxyDelay{}
	route := fmt.Sprintf("/proxies/%s/delay?url=%s&timeout=%d", name, url, timeout)
	err := UnmarshalRequest("get", route, nil, nil, &proxyDelay)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return proxyDelay, nil
}

func SwitchProxy(selector, name string) error {
	route := "/proxies/" + selector
	headers := map[string]string{"Content-Type": "application/json"}
	body := map[string]interface{}{"name": name}

	code, content, err := EasyRequest("put", route, headers, body)
	if err != nil {
		return errors.Trace(err)
	}

	switch code {
	case 204:
		return nil
	case 400, 404:
		return fmt.Errorf(string(content))
	default:
		return fmt.Errorf("unknown error: %s", string(content))
	}
}

type Config struct {
	Port           int      `json:"port"`
	SocksPort      int      `json:"socks-port"`
	RedirPort      int      `json:"redir-port"`
	TproxyPort     int      `json:"tproxy-port"`
	MixedPort      int      `json:"mixed-port"`
	Authentication []string `json:"authentication"`
	AllowLan       bool     `json:"allow-lan"`
	BindAddress    string   `json:"bind-address"`
	Mode           string   `json:"mode"`
	LogLevel       string   `json:"log-level"`
	IPV6           bool     `json:"ipv6"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := UnmarshalRequest("get", "/configs", nil, nil, &config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return config, nil
}

type Rule struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Proxy   string `json:"proxy"`
}

func GetRules() ([]*Rule, error) {
	container := struct {
		Rules []*Rule `json:"rules"`
	}{}
	err := UnmarshalRequest("get", "/rules", nil, nil, &container)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return container.Rules, nil
}

// EnableConfig 这个接口不会影响 external-controller 和 secret 的值
func EnableConfig(path string) error {
	headers := map[string]string{"Content-Type": "application/json"}
	body := map[string]interface{}{"path": path}

	code, content, err := EasyRequest("put", "/configs", headers, body)
	if err != nil {
		return errors.Trace(err)
	}
	if code != 200 {
		return fmt.Errorf("unknown error: %s", string(content))
	}
	return nil
}

func SetConfig(port, socksPort int, redirPort string, allowLan bool, mode, logLevel string) error {
	headers := map[string]string{"Content-Type": "application/json"}
	body := map[string]interface{}{
		"port":       port,
		"socks-port": socksPort,
		"redir-port": redirPort,
		"allow-lan":  allowLan,
		"mode":       mode,
		"log-level":  logLevel,
	}

	code, content, err := EasyRequest("patch", "/configs", headers, body)
	if err != nil {
		return errors.Trace(err)
	}
	if code != 200 {
		return fmt.Errorf("unknown error: %s", string(content))
	}
	return nil
}

// Version Clash 版本信息
type Version struct {
	Meta    bool   `json:"meta"`
	Version string `json:"version"`
}

// GetVersion 获取 Clash 版本信息
func GetVersion() (*Version, error) {
	version := &Version{}
	err := UnmarshalRequest("get", "/version", nil, nil, &version)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return version, nil
}

// Memory 内存占用，单位 kb
//
// 示例: {"inuse":111673344,"oslimit":0}
type Memory struct {
	Inuse   uint64 `json:"inuse"`
	OsLimit uint64 `json:"oslimit"`
}

// GetMemory 获取实时内存占用，单位 kb
func GetMemory(handler func(memory *Memory) (stop bool)) error {
	resp, err := Request("get", "/memory", nil, nil)
	if err != nil {
		return errors.Trace(err)
	}

	HandleStreamResp(resp, func(line []byte) (stop bool) {
		memory := &Memory{}
		if err := json.Unmarshal(line, memory); err != nil {
			return true
		}
		if _stop := handler(memory); _stop {
			return true
		}
		return false
	})
	return nil
}

// Restart 重启内核
func Restart() error {
	code, content, err := EasyRequest("post", "/restart", nil, nil)
	if err != nil {
		return errors.Trace(err)
	}
	if code != 200 {
		return fmt.Errorf("unknown error: %s", string(content))
	}
	return nil
}
