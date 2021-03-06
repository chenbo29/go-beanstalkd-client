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

// init 初始化日志配置
func init() {
	var logPath string
	bsdParamsData = config.GetParams()
	if bsdParamsData.Daemon {
		if os.Getenv("GOOS") == "windows" {
			path, _ := filepath.Abs(os.Args[0])
			logPath = filepath.Dir(path)
		} else {
			logPath = "."
		}
		logFileName = logPath + fmt.Sprintf("/beanstalkf-%s.log", time.Now().Format("2006-01-02"))
		logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			fmt.Println(err)
		}
		logLocal = log.New(logFile, bsdParamsData.Name, log.LstdFlags)
		logLocal.Println(logFileName)
	}
}

// Info info
func Info(v ...interface{}) {
	if bsdParamsData.Daemon {
		logLocal.Println(v)
	} else {
		log.Println(v)
	}
}

// Error error
func Error(v ...interface{}) {
	if bsdParamsData.Daemon {
		logLocal.Println(v)
	} else {
		log.Println(v)
	}
}
