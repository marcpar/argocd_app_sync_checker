package config

import (
	"errors"
	"flag"
)

type ServeFlags struct {
	ConfigPath string
}

func ParseFlags() (*ServeFlags, error) {

	var config_path = flag.String("config", "serve.yaml", "")

	flag.Parse()

	if !flag.Parsed() {
		return nil, errors.New("failed to parse flags")
	}

	return &ServeFlags{
		ConfigPath: *config_path,
	}, nil
}
