package clash

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
)

// RawConfigs (map 指针)存储 GetConfigs() 获取的 Clash 配置
type RawConfigs map[string]interface{}

func GetConfigs() (RawConfigs, error) {
	raw := RawConfigs{}
	code, content, err := EasyRequest("get", "/configs", nil, nil)
	if code != 200 || err != nil {
		return nil, errors.Trace(err)
	}
	err = json.Unmarshal(content, &raw)
	return raw, err
}

func IsTunEnabled(raw RawConfigs) (bool, error) {
	tunField, ok := raw["tun"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("raw has not tun field")
	}
	enable, ok := tunField["enable"].(bool)
	if !ok {
		return false, fmt.Errorf("tun field has not enable field")
	}
	return enable, nil
}

func SetTunEnable(enable bool) error {
	raw := RawConfigs{
		"tun": RawConfigs{
			"enable": enable,
		},
	}
	return SetConfigs(raw)
}

// EnableConfigs 这个接口不会影响 external-controller 和 secret 的值
func EnableConfigs(path string) error {
	headers := map[string]string{"Content-Type": "application/json"}
	body := map[string]interface{}{"path": path}

	code, content, err := EasyRequest("put", "/configs?force=true", headers, body)
	if err != nil {
		return errors.Trace(err)
	}
	if code < 200 || code >= 300 {
		return fmt.Errorf("unknown error: %s", string(content))
	}
	return nil
}

// SetConfigs 更新基本配置，只传入需要修改的配置即可，传入的数据需以 json 格式传入
//
// 命令行调用示例：curl ${controller-api}/configs -X PATCH -d '{"mixed-port": 7890}'
// 参数示例:
//
//	raw := RawConfigs{
//	  "port":       2333,
//	  "socks-port": 2334,
//	  }
func SetConfigs(raw RawConfigs) error {
	headers := map[string]string{"Content-Type": "application/json"}
	code, _, err := EasyRequest("patch", "/configs", headers, raw)
	if err != nil {
		return errors.Trace(err)
	}
	if code < 200 || code >= 300 {
		return fmt.Errorf("return code: %d", code)
	}
	return nil
}

// EnableGeo 更新 GEO 数据库,因更新后会自动重载一次配置
func EnableGeo() error {
	code, content, err := EasyRequest("post", "/configs/geo", nil, nil)
	if err != nil {
		return errors.Trace(err)
	}
	if code < 200 || code >= 300 {
		return fmt.Errorf("unknown error: %s", string(content))
	}
	return nil
}
