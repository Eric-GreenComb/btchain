package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type

func TestZap(t *testing.T) {
	var logI *zap.Logger
	mode := "file"
	env := "debug"

	if mode == "file" {
		var encoder zapcore.Encoder

		if env == "production" {
			encoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
		} else {
			encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		}

		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./log/debug.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		})
		core := zapcore.NewCore(
			encoder,
			w,
			zap.InfoLevel,
		)

		logger := zap.New(core)
		logI = logger

	} else {
		var encoderCfg zapcore.EncoderConfig
		if env == "production" {
			encoderCfg = zap.NewProductionEncoderConfig()
		} else {
			encoderCfg = zap.NewDevelopmentEncoderConfig()
		}

		coreInfo := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(os.Stdout),
			makeInfoFilter(env),
		)
		coreError := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(os.Stderr),
			makeErrorFilter(),
		)

		logI = zap.New(zapcore.NewTee(coreInfo, coreError))
	}

	for i := 0; i < 5; i++ {
		logI.Info("0x000000000000000000000000000000000000000000000000000000000000000000", zap.Int("index", i))
	}

}
