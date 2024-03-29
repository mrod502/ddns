package config

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	CertFilePath   string        `yaml:"cert_file_path"`
	KeyFilePath    string        `yaml:"key_file_path"`
	RemoteHost     string        `yaml:"remote_host"`
	Port           uint16        `yaml:"port"`
	PingInterval   time.Duration `yaml:"ping_interval"`
	PrivateKeyPath string        `yaml:"private_key_path"`
	PublicKeyPath  string        `yaml:"public_key_path"`
	Domain         string        `yaml:"domain"`
	APIKey         string        `yaml:"api_key"`
}

func Parse(path string) (Config, error) {
	var cfg Config
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(b, &cfg)
	fmt.Printf("%+v\n", cfg)
	fmt.Printf("config settings:\n%+v\n", cfg)
	return cfg, err
}
