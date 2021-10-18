package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var log *zap.Logger

func init() {
	//var err error

	// 更改輸出的格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 存入檔案的名稱及位置
	file, err := os.Create("./info.log")
	if err != nil {
		Error("failed to create zap log file, error: " + err.Error())
	}
	writeSyncer := zapcore.AddSync(file)

	// 建立logger
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	log = zap.New(core, zap.AddCallerSkip(1), zap.AddCaller())

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
