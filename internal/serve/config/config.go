package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServeConfig struct {
	ListenAddress  string       `yaml:"listenAddress"`
	ScrapeInterval uint64       `yaml:"scrape_interval"`
	Argocd         ArgocdConfig `yaml:"argocd"`
}

type ArgocdConfig struct {
	Instance string `yaml:"instance"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func ParseConfig(serve_flags ServeFlags) (*ServeConfig, error) {
	config, err := readConfigFromFile(serve_flags.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from path(%s) due to error: %s", serve_flags.ConfigPath, err)
	}

	err = readConfigFromEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func readConfigFromFile(config_path string) (*ServeConfig, error) {

	bytes, err := os.ReadFile(config_path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file(%s) due to error: %s", config_path, err)
	}

	serve_config := ServeConfig{}
	err = yaml.Unmarshal(bytes, &serve_config)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config(%s) due to error: %s", config_path, err)
	}

	return &serve_config, nil
}

// @todo
func readConfigFromEnv(config *ServeConfig) error {
	return nil
}
