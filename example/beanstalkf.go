package main

import (
	"encoding/json"
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
	"github.com/chenbo29/gostorage"

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
			var images []string
			err := json.Unmarshal(body, &images)
			if err != nil {
				logger.Printf("jobID(%d) json decode error %s\n", id, err.Error())
			}
			for _, image := range images {
				uploader := gostorage.NewAliyun()
				result, err := uploader.WebUpload(image)
				if err != nil {
					logger.Printf("jobID(%d) upload error %s\n", id, err.Error())
				} else {
					logger.Printf("jobID(%d) upload success %s\n", id, result.URL)
				}
			}
			return true
		},
	}
	beans.Run(&executeFunc)
}

//func NewLog() *log.Logger {
//	var w io.Writer
//	logFileName := fmt.Sprintf("./beanstalkf-finish-%s.log", time.Now().Format("2006-01-02"))
//	w, _ = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
//	return log.New(w, "go beanstalk client ", log.LstdFlags)
//}
