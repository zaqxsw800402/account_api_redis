package logger_logrus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Logger struct {
	log *logrus.Logger
}

func New(logPath, logFile string) Logger {
	//logFilePath := logPath
	logFileName := logFile
	//fileName := path.Join(logFilePath, logFileName)
	//写入文件
	//src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	src, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err", err)
	}
	//实例化
	logger := logrus.New()
	//设置输出
	logger.Out = src
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//设置日志格式
	//logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return Logger{log: logger}
}

// 日志记录到文件
func (l *Logger) LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		l.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logger) Info(format string) {
	l.log.Info(format)
}

func (l *Logger) Debug(format string) {
	l.log.Debugln(format)
}
