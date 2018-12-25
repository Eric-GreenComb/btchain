package log

import (
	"github.com/axengine/btchain/log"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	Logger = log.Initialize("file", "Initialize", "./log/debug.log", "./log/err.log")
}
