package connect

import (
	"github.com/beanstalkd/go-beanstalk"
	"github.com/chenbo29/go-beanstalkd-client/config"
	"log"
)

// Conn connect beanstalk by params
func Conn(params *config.ParamsData) *beanstalk.Conn {
	beanstalkConn, err := beanstalk.Dial("tcp", params.Host+":"+params.Port)
	if err != nil {
		log.Fatalln(err)
	}
	return beanstalkConn
}
