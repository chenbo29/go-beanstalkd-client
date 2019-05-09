package main

import (
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
)

func main() {
	executeFunc := beans.JobExecuteFunc{
		Execute: func(id uint64, body []byte) bool {
			info := fmt.Sprintf("自定义方法-开始执行jobID[ %d ]jobBody[ %s ]", id, string(body))
			fmt.Println(info)

			if (id % 2) == 0 {
				return false
			} else {
				return true
			}
		},
	}
	beans.Run(&executeFunc)
}
