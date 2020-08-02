package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"app.log"}
	zapLog, err := config.Build(zap.AddStacktrace(zapcore.FatalLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}
	logger = zapLog.Sugar()
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}
