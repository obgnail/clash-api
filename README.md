# clash api

<p align="center">
    <img src="assets/nijika.png" width="450" height="450">
</p>



## 简介

new bing 需要用到香港节点，我平常也不用此节点，于是整了个自动切换脚本。开源 RESTFUL API 部分。



## API

```go
func SetSecrete(sec string) {}
func SetSecreteFromEnv(name string) {}
func SetSecreteFromFile(file string) error {}

func GetLogs(level LogLevel) (chan *Log, error) {}
func GetTraffic(handler func(traffic *Traffic) (stop bool)) error {}
func GetProxies() (map[string]*Proxies, error) {}
func GetProxyMessage(name string) (*Proxy, error) {}
func GetProxyDelay(name string, url string, timeout int) (*ProxyDelay, error) {}
func SwitchProxy(selector, name string) error {}
func GetConfig() (*Config, error) {}
func GetRules() ([]*Rule, error) {}
func EnableConfig(path string) error {}
func SetConfig(port, socksPort int, redirPort string, allowLan bool, mode, logLevel string) error {}
```

 

## Demo

```go
package main

import (
	"fmt"
	"github.com/obgnail/clash-api/clash"
)

func main() {
	if err := clash.SetSecreteFromFile("./secret.txt"); err != nil {
		panic(err)
	}

	proxies, err := clash.GetProxies()
	if err != nil {
		panic(err)
	}
	fmt.Println(proxies)
}
```

