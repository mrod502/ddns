package config

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	cfg := Config{}
	b, _ := yaml.Marshal(cfg)
	fmt.Println(string(b))
}
