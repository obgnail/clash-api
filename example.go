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
