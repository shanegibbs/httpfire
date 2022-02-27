package director

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig() (*Config, error) {
	configFile := os.Getenv("HTTPFIRE_CONFIG")

	configBody, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", configFile, err)
	}

	config := Config{}
	err = yaml.Unmarshal(configBody, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config from file %s: %v", configFile, err)
	}

	return &config, nil
}

type Config struct {
	ListenAddr string     `yaml:"listen_addr"`
	Discovery  *Discovery `yaml:"discovery"`
}

type Discovery struct {
	Static *StaticDiscovery `yaml:"static"`
	DNS    *DNSDiscovery    `yaml:"dns"`
}

type StaticDiscovery struct {
	Endpoints []string `yaml:"endpoints"`
}

type DNSDiscovery struct {
	Name string `yaml:"name"`
	Port uint   `yaml:"port"`
}
