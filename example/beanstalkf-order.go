package main

import (
	"encoding/json"
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
	_ "github.com/chenbo29/gostorage"
	_ "github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	executeFunc := beans.JobExecuteFunc{
		Execute: func(id uint64, body []byte) bool {
			logger := NewLog()
			logger.Printf("jobID(%d) start to do jobBody(%s)", id, string(body))
			var param struct {
				TradeNo int
			}
			err := json.Unmarshal(body, &param)
			if err != nil {
				logger.Printf("error is %s", err)
				return true
			}
			logger.Printf("trade no is %d", param.TradeNo)
			http.Get(fmt.Sprintf("http://video.gutongkj8.com/c?trade_no=%d", param.TradeNo))
			return true
		},
	}
	beans.Run(&executeFunc)
}

func NewLog() *log.Logger {
	var w io.Writer
	logFileName := fmt.Sprintf("./beanstalkf-finish-%s.log", time.Now().Format("2006-01-02"))
	w, _ = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	return log.New(w, "go beanstalk client ", log.LstdFlags)
}
