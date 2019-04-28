package beans

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"go-beanstalkd-client/config"
	"go-beanstalkd-client/connect"
	"os"
)

var con *beanstalk.Conn
var bsdParamsData *config.ParamsData
var separatorLength = 50

func Run()  {
	bsdParamsData = config.GetParams()
	con = connect.Con(bsdParamsData)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "status":
			Status()
		case "start":
			Start()
		case "stop":
			Stop()
		default:
			fmt.Fprintf(os.Stderr, "Usage: %s {start|stop|status}\n", os.Args[0])
		}
	}
	return
}

func Status()  {
	status, _ := con.Stats()
	for key, value := range status {
		fmt.Println(key + GetSeparator(len(key), separatorLength) + value)
	}
	fmt.Println(*bsdParamsData)
}

func Start()  {
	fmt.Println("Start to do something")
	fmt.Println(*bsdParamsData)
}

func Stop()  {
	fmt.Println("Stop Beanstalkd")
	fmt.Println(*bsdParamsData)
}

func GetSeparator(x int, y int) string {
	num := y - x
	separator := " "
	for i:=0; i<num; i++ {
		separator += "-"
	}
	separator += " "
	return separator
}