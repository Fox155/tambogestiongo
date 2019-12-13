package main

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	Sucursal struct {
		Nombre string
	}
	Tiempos struct {
		Dormido string
	}
}

const configFile = "./config.toml"

func initConfig() (*config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	} else if err != nil {
		return nil, err
	}

	conf := &config{}
	if _, err := toml.DecodeFile(configFile, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
