package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
	"io/ioutil"
	"net/http"
)

func main() {
	executeFunc := beans.JobExecuteFunc{
		Execute: func(id uint64, body []byte) bool {
			logger := NewLog()
			logger.Printf("jobID(%d) start to do jobBody(%s)", id, string(body))

			var uri = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
			var params struct {
				AccessToken string
				OpenId      string
			}
			err := json.Unmarshal(body, &params)

			if err != nil {
				logger.Printf("jobID(%d) json decode error %s\n", id, err.Error())
			}

			url := fmt.Sprintf(uri, params.AccessToken, params.OpenId)
			logger.Println(url)

			res, err := http.Get(url)
			if err != nil {
				logger.Printf("jobID(%d) request error %s\n", id, err.Error())
			}
			if res != nil {
				resData, _ := ioutil.ReadAll(res.Body)
				resDataStr := bytes.NewBuffer(resData).String()
				logger.Println(resDataStr)
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
