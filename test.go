package main

import (
	"encoding/json"
	"fmt"
)

type RawConfigs map[string]interface{}

// GetConfigs 获取 Clash 配置, 返回原始 JSON 数据（方便只解析后续部分内容）
func GetConfigs(jsonData []byte) (RawConfigs, error) {
	raw := RawConfigs{}
	err := json.Unmarshal(jsonData, &raw)
	return raw, err
}

func IsTunEnabled(raw RawConfigs) (bool, error) {
	tunValue, ok := raw["tun"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("raw has not tun field")
	}
	enable, ok := tunValue["enable"].(bool)
	if !ok {
		return false, fmt.Errorf("tun field has not enable field")
	}
	return enable, nil
}

func main() {
	// 原始 JSON 数据
	jsonData := []byte(`{
		"port": 8080,
		"socks-port": 8081,
		"mode": "debug",
		"unknownField1": "unknown value 1",
		"unknownField2": "unknown value 2",
		"tun": {
			"enable": false,
			"stack": "system"
		}
	}`)
	raw, _ := GetConfigs(jsonData)
	fmt.Printf("%+v\n", raw)
	fmt.Printf("%+v\n", raw["tun"])
	enabled, _ := IsTunEnabled(raw)
	fmt.Printf("%+v\n", enabled)
}
