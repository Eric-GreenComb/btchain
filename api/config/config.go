package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Bind     string
	RPC      string
	Writable bool
	IsAdmin  bool
	Log      LogInfo
	Cache    Cache
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

type Cache struct {
	RedisStore     bool
	Dial           string
	Password       string
	Db             int
	MaxConnections int
}
