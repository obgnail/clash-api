package main

import (
	"fmt"
	"github.com/obgnail/clash-api/clash"
)

func main() {
	err := clash.SetSecretFromEnv("clash-api-secret")
	checkError(err)

	getProxies()
	getTraffic()
}

func getProxies() {
	proxies, err := clash.GetProxies()
	checkError(err)
	fmt.Println(proxies)
}

func getTraffic() {
	err := clash.GetTraffic(func(traffic *clash.Traffic) (stop bool) {
		fmt.Println(traffic.Up, traffic.Down)
		return false
	})
	checkError(err)

	forever := make(chan struct{}, 1)
	<-forever
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
