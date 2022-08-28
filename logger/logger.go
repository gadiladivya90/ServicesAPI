package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	var err error
	log, err = zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		panic("failed to init Logger!")
	}
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}
func Log(msg string, fields ...zap.Field) {
	log.Log(zapcore.LevelOf(zapcore.InfoLevel), msg)
}
