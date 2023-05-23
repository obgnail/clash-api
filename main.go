package main

import (
	"fmt"
	"github.com/obgnail/clash-api/clash"
)

func testGetTraffic() {
	err := clash.GetTraffic(func(traffic *clash.Traffic) (stop bool) {
		println(traffic.Up, traffic.Down)
		return false
	})
	checkError(err)
}

func testGetLogs() {
	logChan, err := clash.GetLogs(clash.LevelInfo)
	checkError(err)
	for log := range logChan {
		fmt.Println(log)
	}
}

func testGetProxies() {
	proxies, err := clash.GetProxies()
	checkError(err)
	fmt.Println(proxies)
}

func testGetProxyMessage() {
	proxy, err := clash.GetProxyMessage("香港中继 01")
	checkError(err)
	fmt.Println(proxy)
}

func testGetProxyDelay() {
	proxyDelay, err := clash.GetProxyDelay("香港中继 01", "https://bing.com/chat", 3000)
	checkError(err)
	fmt.Println(proxyDelay)
}

func testSelectProxy() {
	err := clash.SwitchProxy("🍃 Proxies", "香港专线 01")
	checkError(err)
	proxies, err := clash.GetProxies()
	checkError(err)
	fmt.Println("proxies", proxies)
}

func testGetConfig() {
	config, err := clash.GetConfig()
	checkError(err)
	fmt.Println("config", config)
}

func testGetRules() {
	rules, err := clash.GetRules()
	checkError(err)
	fmt.Println("rules", rules)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("-----------")
	clash.SetSecreteFromFile("./secret.txt")
	testGetRules()

	//forever := make(chan struct{}, 1)
	//<-forever

	fmt.Println("end")
}
