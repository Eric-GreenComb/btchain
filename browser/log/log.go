package log

import (
	"github.com/axengine/btchain/log"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	Logger = log.Initialize("debug", "Initialize", "debug.log", "err.log")
}
