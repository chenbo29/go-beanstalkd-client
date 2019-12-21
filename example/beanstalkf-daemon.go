package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

const name = "beanstalkf-order"
const pidPath = "./"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [start|stop|status]\n", os.Args[0])
		os.Exit(0)
	}
	switch os.Args[1] {
	case "start":
		Start()
	case "status":
		Start()
	case "stop":
		pid := getPid()
		pro, err := os.FindProcess(pid)
		CheckErr(err)
		CheckErr(pro.Kill())
		fmt.Printf("Pid [%d] Stop \n", pid)
	default:
		fmt.Printf("Usage: %s [start|stop|status]\n", os.Args[0])
		os.Exit(0)
	}

}

// CheckErr 检查error
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// Start 启动
func Start() {
	checkRunningProcess()
	cmdName := getCmdName()
	os.Args[0] = cmdName
	process, err := os.StartProcess(cmdName, os.Args, &os.ProcAttr{Files: []*os.File{nil, os.Stdout, os.Stderr}})
	CheckErr(err)
	fmt.Println("start success")
	recordPid(process.Pid)

	//if os.Stdout == nil && os.Stderr == nil{
	//} else {
	//	fmt.Println("start error")
	//	fmt.Println(os.Stdout)
	//	fmt.Println(os.Stderr)
	//	time.Sleep(3 * time.Second)
	//	stopProcessByPid(process.Pid)
	//}
}

// recordPid 记录进程PID
func recordPid(pid int) {
	var err error
	f, err := os.OpenFile(pidPath+name+".pid", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	CheckErr(err)
	_, err = f.WriteString(strconv.Itoa(pid))
	CheckErr(err)
	CheckErr(f.Close())
}

// getPid 获取进程PID
func getPid() int {
	p, err := ioutil.ReadFile(pidPath + name + ".pid")
	if err != nil {
		return 0
	}
	pid, e := strconv.Atoi(string(p))
	CheckErr(e)
	return pid
}

// checkRunningProcess
func checkRunningProcess() {
	pid := getPid()
	if pid != 0 {
		stopProcessByPid(pid)
	}
}

// stopProcessByPid stop process by pid
func stopProcessByPid(pid int) {
	pro, err := os.FindProcess(pid)
	if err == nil {
		err := pro.Kill()
		CheckErr(err)
	}
}

func getCmdName() string {
	var cmdName string
	if os.Getenv("GOEXE") == ".exe" {
		cmdName = name + os.Getenv("GOEXE")
	} else {
		cmdName = name
	}
	lp, err := exec.LookPath(cmdName)
	CheckErr(err)
	return lp
}
