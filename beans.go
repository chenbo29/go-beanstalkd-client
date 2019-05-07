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
	"github.com/beanstalkd/go-beanstalk"
	"github.com/chenbo29/go-beanstalkd-client/config"
	"github.com/chenbo29/go-beanstalkd-client/connect"
	"github.com/chenbo29/go-beanstalkd-client/loglocal"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

var commentMap = map[string]string{"cmd-put": "总共执行put指令的次数", "cmd-peek": "总共执行peek指令的次数", "cmd-peek-ready": "总共执行peek-ready指令的次数", "cmd-peek-delayed": "总共执行peek-delayed指令的次数", "cmd-peek-buried": "总共执行peek-buried指令的次数", "cmd-reserve": "总共执行reserve指令的次数", "cmd-use": "总共执行use指令的次数", "cmd-watch": "总共执行watch指令的次数", "cmd-ignore": "总共执行ignore指令的次数", "cmd-release": "总共执行release指令的次数", "cmd-bury": "总共执行bury指令的次数", "cmd-kick": "总共执行kick指令的次数", "cmd-stats": "总共执行stats指令的次数", "cmd-stats-job": "总共执行stats-job指令的次数", "cmd-stats-tube": "总共执行stats-tube指令的次数", "cmd-list-tubes": "总共执行list-tubes指令的次数", "cmd-list-tube-used": "总共执行list-tube-used指令的次数", "cmd-list-butes-watched": "总共执行list-tubes-watched指令的次数", "cmd-pause-tube": "总共执行pause-tube指令的次数", "job-timeouts": "所有超时的job的总共数量", "max-job-size": "job的数据部分最大长度", "current-tubes": "当前存在的tube数量", "current-connections": "当前打开的连接数", "current-producers": "当前所有的打开的连接中至少执行一次put指令的连接数量", "current-workers": "当前所有的打开的连接中至少执行一次reserve指令的连接数量", "current-waiting": "当前所有的打开的连接中执行reserve指令但是未响应的连接数量", "total-connections": "总共处理的连接数", "pid": "服务器进程的id", "version": "服务器版本号", "rusage-utime": "进程总共占用的用户CPU时间", "rusage-stime": "进程总共占用的系统CPU时间", "uptime": "服务器进程运行的秒数", "binlog-oldest-index": "开始储存jobs的binlog索引号", "binlog-current-index": "当前储存jobs的binlog索引号", "binlog-max-size": "binlog的最大容量", "binlog-records-written": "binlog累积写入的记录数", "binlog-records-migrated": "is the cumulative number of records written as part of compaction.", "id": "一个随机字符串，在beanstalkd进程启动时产生", "hostname": "主机名", "name": "表示tube的名称", "current-jobs-urgent": "此tube中优先级小于1024状态为ready的job数量", "current-jobs-ready": "此tube中状态为ready的job数量", "current-jobs-reserved": "此tube中状态为reserved的job数量", "current-jobs-delayed": "此tube中状态为delayed的job数量", "current-jobs-buried": "此tube中状态为buried的job数量", "total-jobs": "此tube中创建的所有job数量", "current-using": "使用此tube打开的连接数", "current-wating": "使用此tube打开连接并且等待响应的连接数", "current-watching": "打开的连接监控此tube的数量", "pause": "此tube暂停的秒数", "cmd-delete": "此tube中总共执行的delete指令的次数", "pause-time-left": "此tube暂停剩余的秒数"}
var conn *beanstalk.Conn
var bsdParamsData *config.ParamsData
var separatorLength = 50
var tubesChan chan string
var commandName = "go-beanstalk-client"
var Daemon = false

const workerNum = 10
const reserveTime = 5

func Run() {
	bsdParamsData = config.GetParams()
	conn = connect.Conn(bsdParamsData)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "status":
			Status()
		case "start":
			Start()
		case "testPut":
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
			tube := beanstalk.Tube{Conn: conn, Name: tubeName}
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

func ListTubeInfo(t *beanstalk.Tube) {
	info, err := t.Stats()
	if err != nil {
		log.Println(err)
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
		tube := beanstalk.Tube{Conn: conn, Name: tubeName}
		info, _ := tube.Stats()
		loglocal.Info(info)
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
		loglocal.Info(key + GetSeparator(len(key), separatorLength) + (*status)[key] + " [" + commentMap[key] + "]")
	}
}

func TestPut(tubeName *string) {
	tube := beanstalk.Tube{Conn: conn, Name: *tubeName}
	for i := 0; i < 100; i++ {
		info := []byte(*tubeName + " test info " + strconv.Itoa(i))
		jobId, _ := tube.Put(info, 0, 0, 10*time.Second)
		fmt.Println(jobId)
	}
	ListTubeInfo(&tube)
}

func TubeFactoryStart(tubeName string) {
	paramsData := config.GetParams()
	conn := connect.Conn(paramsData)
	tf := NewTubeFactory(tubeName, workerNum, conn)
	tf.Run()
}

// 厂长
func Monitor(originTubeNum int) {
	for {
		TubesName, _ := conn.ListTubes()
		TubeNum := len(TubesName)
		if TubeNum > originTubeNum {
			for x := originTubeNum; x < TubeNum; x++ {
				loglocal.Info(fmt.Sprintf("Monitor TubeFactory(%s) Start", TubesName[x]))
				go TubeFactoryStart(TubesName[x])
			}
			originTubeNum = TubeNum
		}
	}
}

// Tube工厂
type TubeFactory struct {
	workerNum int
	name      string
	conn      *beanstalk.Conn
}

// 创建Tube工厂
func NewTubeFactory(name string, num int, conn *beanstalk.Conn) *TubeFactory {
	w := TubeFactory{
		workerNum: num,
		name:      name,
		conn:      conn,
	}
	return &w
}

// 工厂启动
func (tf *TubeFactory) Run() {
	loglocal.Info(fmt.Sprintf("TubeFactory(%s) Running, %d`s Worker", tf.name, tf.workerNum))
	for i := 0; i < tf.workerNum; i++ {
		w := NewWorker(strconv.Itoa(i), func(name string, conn *beanstalk.Conn) error {
			tubeSet := beanstalk.NewTubeSet(conn, tf.name)
			var errorDeleteJob error
			var errorReserve error
			//var errorBury error
			var jobId uint64
			var jobBody []byte
			for {
				jobId, jobBody, errorReserve = tubeSet.Reserve(reserveTime)
				if errorReserve != nil {
					//loglocal.Error(fmt.Sprintf("%s Error: %s", tf.name, errorReserve))
				} else {
					loglocal.Info(fmt.Sprintf("%s Worker(%s) Get JobId [%d] JobBody [%s] And Start To Do", tf.name, name, jobId, string(jobBody)))

					for {
						if errorDeleteJob = tf.conn.Delete(jobId); errorDeleteJob != nil {
							time.Sleep(1 * time.Second)
							//loglocal.Error(errorDeleteJob)
							//loglocal.Error(tf.conn.StatsJob(jobId))
							continue
						} else {
							loglocal.Info(fmt.Sprintf("%s Worker(%s) Start To Do Job(%d) Finish !", tf.name, name, jobId))
							break
						}
					}
				}
				time.Sleep(time.Second * 1)
			}
			return nil
		})
		go w.Execute(tf)
	}
}

// 工人
type Worker struct {
	name     string                                        // 名称
	f        func(name string, conn *beanstalk.Conn) error //操作步骤内容
	conn     *beanstalk.Conn
	tubeName string
}

// 分配工人
func NewWorker(name string, f func(name string, conn *beanstalk.Conn) error) *Worker {
	w := Worker{
		name: name,
		f:    f,
	}
	return &w
}

// 工人开始操作
func (w *Worker) Execute(tf *TubeFactory) {
	bsdParamsData = config.GetParams()
	conn = connect.Conn(bsdParamsData)
	_ = w.f(w.name, conn)
}

// 获取任务Job
func (w *Worker) ReserveJob() {
	tubeSet := beanstalk.NewTubeSet(w.conn, w.tubeName)
	for {
		jobId, jobBody, err := tubeSet.Reserve(reserveTime)
		if err != nil {
			loglocal.Error(fmt.Sprintf("%s Error: %s", w.tubeName, err))
		} else {
			loglocal.Info(fmt.Sprintf("%s Get JobId [%d] JobBody [%s]", w.tubeName, jobId, string(jobBody)))
		}
		time.Sleep(5 * time.Second)
	}
}
