package main

import (
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	executeFunc := beans.JobExecuteFunc{
		Execute: func(id uint64, body []byte) bool {
			log := NewLog()
			log.Println("自定义方法-开始执行jobID[ %d ]jobBody[ %s ]", id, string(body))
			if (id % 2) == 0 {
				return false
			} else {
				return true
			}
		},
	}
	beans.Run(&executeFunc)
}

func NewLog() *log.Logger {
	var w io.Writer
	logFileName := fmt.Sprintf("./beanstalkf-finish-%s.log", time.Now().Format("2006-01-02"))
	w, _ = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	return log.New(w, "result", log.LstdFlags)
}
