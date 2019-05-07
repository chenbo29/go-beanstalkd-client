package beans

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/chenbo29/go-beanstalkd-client/config"
	"github.com/chenbo29/go-beanstalkd-client/connect"
	"github.com/chenbo29/go-beanstalkd-client/loglocal"
	"time"
)

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
