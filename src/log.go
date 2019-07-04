package src

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger *logrus.Logger

func InitLogger() {
	fmt.Println("init log")
	logFileName := time.Now().Format("./log/20060102-150405.log")
	logFile, err := createFile(logFileName)
	if nil != err {
		fmt.Printf("日志文件[%s]打开失败\n", logFileName)
		return
	}
	Logger = &logrus.Logger{
		Out: logFile,
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
		ExitFunc: func(i int) {
			if nil != logFile {
				logFile.Close()
			}
			os.Exit(i)
		},
		Level: logrus.InfoLevel,
	}
}