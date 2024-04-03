package clash

import (
	"fmt"
	"github.com/juju/errors"
	"io/ioutil"
	"os"
)

var Secret = ""

func SetSecret(secret string) {
	Secret = secret
}

func SetSecretFromEnv(name string) error {
	secrete := os.Getenv(name)
	if len(secrete) != 0 {
		Secret = secrete
		return nil
	}
	return fmt.Errorf("has no such name")
}

func SetSecretFromFile(file string) error {
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return errors.Trace(err)
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Trace(err)
	}
	Secret = string(content)
	return nil
}
