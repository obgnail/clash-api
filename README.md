# clash api

<p align="center">
    <img src="assets/nijika.png" width="450" height="450">
</p>
本项目是对 [clash doc](https://clash.gitbook.io/doc/) 的翻译。我丢，搬砖活。




## Introduction

> Clash RESTful API 是一套控制 Clash 的一个途径，能获取 Clash 中的一些信息，同时也能控制 Clash 内部的配置。基于 API，可以打造自己的可视化操作部分，也是实现 Clash GUI 的重要组成部分。—— Clash doc

因为 new bing 需要用到香港节点，我平常也不用此节点，于是整了个自动切换脚本。开源 RESTful API 部分。



## API

```go
// 设置密钥(调用业务接口前必须调用此接口)
func SetSecrete(sec string) {}
func SetSecreteFromEnv(name string) {}
func SetSecreteFromFile(file string) error {}

// 监控请求日志
func GetLogs(level LogLevel) (chan *Log, error) {}

// 每秒推送一次，上下载流量
func GetTraffic(handler func(traffic *Traffic) (stop bool)) error {}

// 节点列表
func GetProxies() (map[string]*Proxies, error) {}

// 具体节点的信息
func GetProxyMessage(name string) (*Proxy, error) {}

// 测试具体节点的延时
func GetProxyDelay(name string, url string, timeout int) (*ProxyDelay, error) {}

// 切换节点
func SwitchProxy(selector, name string) error {}

// 获取配置
func GetConfig() (*Config, error) {}

// 设置配置
func SetConfig(port, socksPort int, redirPort string, allowLan bool, mode, logLevel string) error {}

// 应用配置(不会影响 external-controller 和 secret 的值)
func EnableConfig(path string) error {}

// PAC rules
func GetRules() ([]*Rule, error) {}
```

 

## Example

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

