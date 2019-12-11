package main

import (
	"encoding/json"
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
	_ "github.com/chenbo29/gostorage"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	executeFunc := beans.JobExecuteFunc{
		Execute: func(id uint64, body []byte) bool {
			logger := NewLog()
			logger.Printf("jobID(%d) start to do jobBody(%s)", id, string(body))
			var param struct {
				FollowerNum      int
				InfluencerUserId int
				InfluencerPFId   int
			}
			err := json.Unmarshal(body, &param)
			if err != nil {
				log.Fatal(err)
			}
			logger.Printf("follower_num is %d", param.FollowerNum)
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
