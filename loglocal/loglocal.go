package loglocal

import (
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client/config"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logFileName string
var logLocal *log.Logger
var bsdParamsData *config.ParamsData

// 初始化日志配置
func init() {
	bsdParamsData = config.GetParams()
	if bsdParamsData.Daemon {
		path, _ := filepath.Abs(os.Args[0])
		logPath := filepath.Dir(path)
		logFileName = logPath + fmt.Sprintf("\\%s.log", time.Now().Format("2006-01-02"))
		fmt.Println("Log Init")
		fmt.Println(logFileName)
		logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			panic(err)
		}
		logLocal = log.New(logFile, bsdParamsData.Name, log.LstdFlags)
	}
}

// 记录信息
func Info(v ...interface{}) {
	if bsdParamsData.Daemon {
		logLocal.Println(v)
	} else {
		log.Println(v)
	}
}

// 记录错误
func Error(v ...interface{}) {
	if bsdParamsData.Daemon {
		logLocal.Println(v)
	} else {
		log.Println(v)
	}
}
