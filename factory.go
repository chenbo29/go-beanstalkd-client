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
	workerNum int
	name      string
	conn      *beanstalk.Conn
}

// NewTubeFactory 创建Tube工厂
func NewTubeFactory(name string, num int, conn *beanstalk.Conn) *TubeFactory {
	w := TubeFactory{
		workerNum: num,
		name:      name,
		conn:      conn,
	}
	return &w
}

// Run 工厂启动
func (tf *TubeFactory) Run() {
	loglocal.Info(fmt.Sprintf("TubeFactory(%s) Running, %d`s Worker", tf.name, tf.workerNum))
	for i := 0; i < tf.workerNum; i++ {
		w := NewWorker(strconv.Itoa(i), func(name string, conn *beanstalk.Conn) error {
			tubeSet := beanstalk.NewTubeSet(conn, tf.name)
			var errorDeleteJob error
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

					for {
						if errorDeleteJob = tf.conn.Delete(jobID); errorDeleteJob != nil {
							time.Sleep(1 * time.Second)
							//loglocal.Error(errorDeleteJob)
							//loglocal.Error(tf.conn.StatsJob(jobId))
							continue
						} else {
							loglocal.Info(fmt.Sprintf("%s Worker(%s) Start To Do Job(%d) Finish !", tf.name, name, jobID))
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
