package clash

import (
	"testing"
)

func TestGetConfigs(t *testing.T) {
	Init()
	configs, err := GetConfigs()
	if err != nil {
		t.Errorf("Error retrieving configs: %s", err)
		return
	}
	t.Logf("Configs: %+v", configs)
}

func TestSetTunEnable(t *testing.T) {
	Init()
	raw, _ := GetConfigs()
	enable, _ := IsTunEnabled(raw)
	t.Logf("Tun enabled: %t", enable)

	err := SetTunEnable(!enable)
	if err != nil {
		t.Errorf("Error switch tun: %s", err)
		return
	}
}
