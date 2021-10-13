package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var log *zap.Logger

func init() {
	var err error

	//config := zap.NewProductionConfig()

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	//config.EncoderConfig = encoderConfig
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	//encoder := zapcore.NewConsoleEncoder(encoderConfig)

	file, _ := os.Create("./info.log")
	writeSyncer := zapcore.AddSync(file)

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	log = zap.New(core, zap.AddCallerSkip(1), zap.AddCaller())
	//log, err = config.Build(zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}
