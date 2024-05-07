package config

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func ConfigureZapLogger(folderName string) *zap.Logger {
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf(`logs/%s/harmonium.log`, folderName),
		MaxSize:    1024, // MB
		MaxBackups: 10,
		MaxAge:     60, // days
		Compress:   true,
	})
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), ws, zap.NewAtomicLevelAt(zap.InfoLevel))
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}
