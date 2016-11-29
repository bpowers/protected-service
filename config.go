// Copyright 2016 Bobby Powers. All rights reserved.

package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type ServiceConfig struct {
	SessionKey string
	CookieName string
	AuthDir    string
}

func ReadConfig(path string) (*ServiceConfig, error) {
	config := new(ServiceConfig)
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		return nil, fmt.Errorf("toml.DecodeFile(%s): %s", path, err)
	}
	return config, nil
}
