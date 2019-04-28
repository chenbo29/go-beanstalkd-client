package beans

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/chenbo29/go-beanstalkd-client/config"
	"github.com/chenbo29/go-beanstalkd-client/connect"
	"os"
	"sort"
)

var conn *beanstalk.Conn
var bsdParamsData *config.ParamsData
var separatorLength = 50

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
		default:
			fmt.Fprintf(os.Stderr, "Usage: %s {start|stop|status}\n", os.Args[0])
		}
	}
	return
}

func Status() {
	tubeName := bsdParamsData.Tube
	if tubeName != "" {
		if tubeName == "all" {
			ListTubesInfo()
		} else {
			tube := beanstalk.Tube{conn, tubeName}
			ListTubeInfo(&tube)
		}
	} else {
		//todo 打印服务状态信息参数
		status, _ := conn.Stats()
		for key, value := range status {
			fmt.Println(key + GetSeparator(len(key), separatorLength) + value)
		}
		tubes, _ := conn.ListTubes()
		fmt.Println(len(tubes))
		fmt.Println(*bsdParamsData)
	}
}

func Start() {
	currentTubesName, _ := conn.ListTubes()
	tubes := make([]beanstalk.Tube, 0)
	for _, currentTubeName := range currentTubesName {
		fmt.Println("current tube:", currentTubeName)
		tubes = append(tubes, beanstalk.Tube{conn, currentTubeName})
	}
	for _, tube := range tubes {
		ListTubeInfo(&tube)
		fmt.Println("")
	}
}

func Stop() {
	fmt.Println("Stop Beanstalkd")
	fmt.Println(*bsdParamsData)
}

func ListTubeInfo(t *beanstalk.Tube) {
	info, _ := t.Stats()
	info["Tube`s Status Info"] = fmt.Sprintf("[%s]", t.Name)
	temp := make([]string, len(info))
	i := 0
	for k := range info {
		temp[i] = k
		i++
	}
	sort.Strings(temp)
	for _, key := range temp {
		fmt.Println(key + GetSeparator(len(key), separatorLength) + info[key])
	}
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
