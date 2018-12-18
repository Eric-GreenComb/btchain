package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Bind         string
	RPC          string
	AdminAddress string
	AdminPrivKey string
	Log          LogInfo
}

func New() *Config {
	return &Config{}
}

func (p *Config) Init(cfgFile string) error {
	_, err := toml.DecodeFile(cfgFile, p)
	return err
}

type LogInfo struct {
	Path string
}
