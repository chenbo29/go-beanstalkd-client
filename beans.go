// Copyright 2018 chenbo29
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package beans implements a queue of beanstalk framework
package beans

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/chenbo29/go-beanstalkd-client/config"
	"github.com/chenbo29/go-beanstalkd-client/connect"
	"os"
	"sort"
	"time"
)

var commentMap = map[string]string{"cmd-put": "总共执行put指令的次数", "cmd-peek": "总共执行peek指令的次数", "cmd-peek-ready": "总共执行peek-ready指令的次数", "cmd-peek-delayed": "总共执行peek-delayed指令的次数", "cmd-peek-buried": "总共执行peek-buried指令的次数", "cmd-reserve": "总共执行reserve指令的次数", "cmd-use": "总共执行use指令的次数", "cmd-watch": "总共执行watch指令的次数", "cmd-ignore": "总共执行ignore指令的次数", "cmd-release": "总共执行release指令的次数", "cmd-bury": "总共执行bury指令的次数", "cmd-kick": "总共执行kick指令的次数", "cmd-stats": "总共执行stats指令的次数", "cmd-stats-job": "总共执行stats-job指令的次数", "cmd-stats-tube": "总共执行stats-tube指令的次数", "cmd-list-tubes": "总共执行list-tubes指令的次数", "cmd-list-tube-used": "总共执行list-tube-used指令的次数", "cmd-list-butes-watched": "总共执行list-tubes-watched指令的次数", "cmd-pause-tube": "总共执行pause-tube指令的次数", "job-timeouts": "所有超时的job的总共数量", "max-job-size": "job的数据部分最大长度", "current-tubes": "当前存在的tube数量", "current-connections": "当前打开的连接数", "current-producers": "当前所有的打开的连接中至少执行一次put指令的连接数量", "current-workers": "当前所有的打开的连接中至少执行一次reserve指令的连接数量", "current-waiting": "当前所有的打开的连接中执行reserve指令但是未响应的连接数量", "total-connections": "总共处理的连接数", "pid": "服务器进程的id", "version": "服务器版本号", "rusage-utime": "进程总共占用的用户CPU时间", "rusage-stime": "进程总共占用的系统CPU时间", "uptime": "服务器进程运行的秒数", "binlog-oldest-index": "开始储存jobs的binlog索引号", "binlog-current-index": "当前储存jobs的binlog索引号", "binlog-max-size": "binlog的最大容量", "binlog-records-written": "binlog累积写入的记录数", "binlog-records-migrated": "is the cumulative number of records written as part of compaction.", "id": "一个随机字符串，在beanstalkd进程启动时产生", "hostname": "主机名", "name": "表示tube的名称", "current-jobs-urgent": "此tube中优先级小于1024状态为ready的job数量", "current-jobs-ready": "此tube中状态为ready的job数量", "current-jobs-reserved": "此tube中状态为reserved的job数量", "current-jobs-delayed": "此tube中状态为delayed的job数量", "current-jobs-buried": "此tube中状态为buried的job数量", "total-jobs": "此tube中创建的所有job数量", "current-using": "使用此tube打开的连接数", "current-wating": "使用此tube打开连接并且等待响应的连接数", "current-watching": "打开的连接监控此tube的数量", "pause": "此tube暂停的秒数", "cmd-delete": "此tube中总共执行的delete指令的次数", "pause-time-left": "此tube暂停剩余的秒数"}
var conn *beanstalk.Conn
var bsdParamsData *config.ParamsData
var separatorLength = 50
var tubesChan chan string

func Run() {
	bsdParamsData = config.GetParams()
	conn = connect.Conn(bsdParamsData)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "status":
			Status()
		case "start":
			Start()
		case "stop":
			Stop()
		case "test":
			TestPut(&os.Args[2])
		default:
			fmt.Fprintf(os.Stderr, "Usage: %s {start|stop|status}\n", os.Args[0])
		}
	}
	return
}

func Status() {
	if bsdParamsData.IsAllStatus {
		statusMap, _ := conn.Stats()
		statusSlice := GetSliceByMapString(statusMap)
		ShowStatus(&statusSlice, &statusMap)
	} else {
		tubeName := bsdParamsData.Tube
		if tubeName == "all" {
			ListTubesInfo()
		} else {
			tube := beanstalk.Tube{conn, tubeName}
			ListTubeInfo(&tube)
		}
	}
}

// 获取tube的状态信息
func Start() {
	go Monitor(0)
	for {
		time.Sleep(1 * time.Second)
	}
}

func Stop() {
	procAttr := &os.ProcAttr{
		Dir:   "D:/code/go-project/",
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	process, err := os.StartProcess("chenbo", os.Args, procAttr)
	if err != nil {
		panic(err)
	}
	fmt.Println(process)
	fmt.Println("Stop Beanstalkd")
}

func ListTubeInfo(t *beanstalk.Tube) {
	info, err := t.Stats()
	if err != nil {
		logs.Error(err)
		return
	}
	info["Tube`s Status Info"] = fmt.Sprintf("[%s]", t.Name)
	infoSlice := GetSliceByMapString(info)
	ShowStatus(&infoSlice, &info)
	return
}

func ListTubesInfo() {
	tubesName, _ := conn.ListTubes()
	for _, tubeName := range tubesName {
		tube := beanstalk.Tube{conn, tubeName}
		info, _ := tube.Stats()
		fmt.Println(info)
	}
}

func GetSeparator(x int, y int) string {
	num := y - x
	separator := " "
	for i := 0; i < num; i++ {
		separator += "-"
	}
	separator += " "
	return separator
}

func GetSliceByMapString(m map[string]string) []string {
	temp := make([]string, len(m))
	i := 0
	for key := range m {
		temp[i] = key
		i++
	}
	sort.Strings(temp)
	return temp
}

func ShowStatus(s *[]string, status *map[string]string) {
	for _, key := range *s {
		fmt.Println(key + GetSeparator(len(key), separatorLength) + (*status)[key] + " [" + commentMap[key] + "]")
	}
}

func TestPut(tubeName *string) {
	tube := beanstalk.Tube{conn, *tubeName}
	info := []byte(*tubeName + " test info")
	jobId, _ := tube.Put(info, 0, 0, 10)
	fmt.Println(jobId)

	ListTubeInfo(&tube)
}

func Work(tubeName *string) {
	workConn := connect.Conn(bsdParamsData)
	for {
		tubeSet := beanstalk.NewTubeSet(workConn, *tubeName)
		jobId, jobBody, err := tubeSet.Reserve(5 * time.Second)
		if err != nil {
			errorInfo := fmt.Sprintf("%s [%s]", err, *tubeName)
			logs.Error(errorInfo)
			continue
		}
		info := fmt.Sprintf("Tube[%s] JobId[%d] JobBody[%s]", *tubeName, jobId, string(jobBody))
		logs.Info(info)
		// todo 处理队列任务
		workConn.Delete(jobId)
	}
}

func Monitor(originTubeNum int) {
	for true {
		TubesName, _ := conn.ListTubes()
		TubeNum := len(TubesName)
		if TubeNum > originTubeNum {
			for x := originTubeNum; x < TubeNum; x++ {
				logs.Info(fmt.Sprintf("Monitor Tube [%s] Reserve Worker Start", TubesName[x]))
				go Work(&TubesName[x])
			}
			originTubeNum = TubeNum
		}
	}
}
