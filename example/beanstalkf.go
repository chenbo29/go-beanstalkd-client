package main

import (
	"encoding/json"
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
	"github.com/chenbo29/gostorage"

	//"github.com/chenbo29/gostorage"
	_ "github.com/chenbo29/gostorage"
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

			var images []string
			json.Unmarshal(body, &images)
			fmt.Println(images)
			for _, image := range images {
				//fmt.Println(image)
				uploader := gostorage.NewAliyun()
				result, _ := uploader.WebUpload(image)
				fmt.Println(result)
			}
			return true
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
