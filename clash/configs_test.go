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
