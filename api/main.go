package main

import (
	"github.com/axengine/btchain/api/config"
	"github.com/axengine/btchain/api/server"
	"github.com/axengine/btchain/log"
	"go.uber.org/zap"
	"path"
)

func main() {
	var logger *zap.Logger
	cfg := config.New()
	if err := cfg.Init("./config/config.toml"); err != nil {
		panic("On init yaml:" + err.Error())
	}

	logger = log.Initialize("file", "debug", path.Join(cfg.Log.Path, "api.debug.log"))
	logger.Info("config", zap.Any("cfg", cfg))
	server := server.NewServer(logger, cfg)
	server.Start()
}
