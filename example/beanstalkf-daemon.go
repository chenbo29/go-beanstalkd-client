package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

const name = "beanstalkf"
const pidPath = "/var/run/"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [start|stop|status]\n", os.Args[0])
		os.Exit(1)
	}
	switch os.Args[1] {
	case "start":
		Start()
	case "status":
		Start()
	case "stop":
		pid := GetPid()
		pro, err := os.FindProcess(pid)
		CheckErr(err)
		CheckErr(pro.Kill())
		fmt.Printf("Pid [%d] Stop \n", pid)
	default:
		fmt.Printf("Usage: %s [start|stop|status]\n", os.Args[0])
		os.Exit(0)
	}

}

// RecordPid 记录进程PID
func RecordPid(pro *os.Process) {
	var err error
	f, err := os.OpenFile(pidPath+name+".pid", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	CheckErr(err)
	_, err = f.WriteString(strconv.Itoa(pro.Pid))
	CheckErr(err)
	CheckErr(f.Close())
}

// GetPid 获取进程PID
func GetPid() int {
	p, err := ioutil.ReadFile(pidPath + name + ".pid")
	CheckErr(err)
	pid, e := strconv.Atoi(string(p))
	CheckErr(e)
	return pid
}

// CheckErr 检查error
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Start 启动
func Start() {
	exe := os.Getenv("GOEXE")
	var cmdName string
	if exe == ".exe" {
		cmdName = name + os.Getenv("GOEXE")
	} else {
		cmdName = name
	}
	lp, err := exec.LookPath(cmdName)
	CheckErr(err)
	cmdName = lp
	procAttr := &os.ProcAttr{
		Files: []*os.File{nil, nil, os.Stderr},
	}

	args := []string{cmdName}
	for key, value := range os.Args {
		if key == 0 {
			continue
		} else {
			args = append(args, value)
			if value == "start" {
				args = append(args, "-d=true")
			}
		}
	}
	process, err := os.StartProcess(cmdName, args, procAttr)
	CheckErr(err)
	RecordPid(process)
}
