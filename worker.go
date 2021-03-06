package beans

import (
	"github.com/beanstalkd/go-beanstalk"
)

// Worker 工人
type Worker struct {
	name     string                                        // 名称
	f        func(name string, conn *beanstalk.Conn) error //操作步骤内容
	conn     *beanstalk.Conn
	tubeName string
}

// NewWorker 分配工人
func NewWorker(name string, f func(name string, conn *beanstalk.Conn) error) *Worker {
	w := Worker{
		name: name,
		f:    f,
	}
	return &w
}

// Execute 工人开始操作
func (w *Worker) Execute(tf *TubeFactory) {
	//bsdParamsData = config.GetParams()
	//conn = connect.Conn(bsdParamsData)
	_ = w.f(w.name, tf.conn)
}
