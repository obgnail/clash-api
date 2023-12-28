package main

import (
	"fmt"
	"github.com/obgnail/clash-api/clash"
)

func main() {
	if err := clash.SetSecretFromEnv("clash-api-secret"); err != nil {
		panic(err)
	}

	proxies, err := clash.GetProxies()
	if err != nil {
		panic(err)
	}
	fmt.Println(proxies)
}
