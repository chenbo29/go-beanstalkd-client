package connect

import (
	"github.com/beanstalkd/go-beanstalk"
	"go-beanstalkd-client/config"
)

func Con(params *config.ParamsData) *beanstalk.Conn {
	con, err := beanstalk.Dial("tcp", params.Host + ":" + params.Port)
	if err != nil {
		panic(err)
	}
	return con
}