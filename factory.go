package beans

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/chenbo29/go-beanstalkd-client/loglocal"
	"strconv"
	"time"
)

// TubeFactory Tube工厂
type TubeFactory struct {
	workerNum   int
	name        string
	conn        *beanstalk.Conn
	executeFunc *JobExecuteFunc
}

// NewTubeFactory 创建Tube工厂
func NewTubeFactory(name string, num int, conn *beanstalk.Conn, executeFunc *JobExecuteFunc) *TubeFactory {
	w := TubeFactory{
		workerNum:   num,
		name:        name,
		conn:        conn,
		executeFunc: executeFunc,
	}
	return &w
}

// Run 工厂启动
func (tf *TubeFactory) Run() {
	loglocal.Info(fmt.Sprintf("TubeFactory(%s) Running, %d`s Worker", tf.name, tf.workerNum))
	for i := 0; i < tf.workerNum; i++ {
		w := NewWorker(strconv.Itoa(i), func(name string, conn *beanstalk.Conn) error {
			executeFunc := *tf.executeFunc

			tubeSet := beanstalk.NewTubeSet(tf.conn, tf.name)
			var errorReserve error
			//var errorBury error
			var jobID uint64
			var jobBody []byte
			for {
				jobID, jobBody, errorReserve = tubeSet.Reserve(reserveTime)
				if errorReserve != nil {
					//loglocal.Error(fmt.Sprintf("%s Error: %s", tf.name, errorReserve))
				} else {
					loglocal.Info(fmt.Sprintf("%s Worker(%s) Get JobId [%d] JobBody [%s] And Start To Do", tf.name, name, jobID, string(jobBody)))
					if executeFunc.Execute(jobID, jobBody) {
						// 业务执行结果成功
						DeleteJob(tf, name, jobID)
					} else {
						// 业务执行结果失败
						BuryJob(tf, name, jobID)
					}
				}
				time.Sleep(time.Second * 1)
			}
			return errorReserve
		})
		go w.Execute(tf)
	}
}

// DeleteJob 业务函数执行成功后删除job
func DeleteJob(tf *TubeFactory, workerName string, jobID uint64) {
	var errorDeleteJob error
	for {
		if errorDeleteJob = tf.conn.Delete(jobID); errorDeleteJob != nil {
			time.Sleep(1 * time.Second)
			//loglocal.Error(errorDeleteJob)
			//loglocal.Error(tf.conn.StatsJob(jobId))
			continue
		} else {
			loglocal.Info(fmt.Sprintf("%s Worker(%s) Start To Do Job(%d) Finish ✔ !", tf.name, workerName, jobID))
			break
		}
	}
}

// BuryJob 回收Job
func BuryJob(tf *TubeFactory, workerName string, jobID uint64) {
	if err := tf.conn.Bury(jobID, 0); err != nil {
		loglocal.Info(fmt.Sprintf("%s Worker(%s) Start To Do Job(%d) Failed ⚠ !", tf.name, workerName, jobID))
		loglocal.Error(err)
	} else {
		loglocal.Info(fmt.Sprintf("%s Worker(%s) Start To Do Job(%d) Failed And Buried ⚠ !", tf.name, workerName, jobID))
	}
}
