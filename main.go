package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [command]\n", os.Args[0])
		os.Exit(1)
	}
	switch os.Args[1] {
	case "start":
		Start()
	case "stop":
		pid := GetPid()
		pro, err := os.FindProcess(pid)
		CheckErr(err)
		CheckErr(pro.Kill())
		fmt.Printf("Pid [%d] Stop \n", pid)
	default:
		fmt.Printf("Usage: %s [start|stop]\n", os.Args[0])
		os.Exit(1)
	}

}

func RecordPid(pro *os.Process)  {
	fmt.Println("Record Pid", pro.Pid)
	f,_ := os.OpenFile("beans.pid", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	_, err := f.WriteString(strconv.Itoa(pro.Pid))
	if err != nil {
		panic(err)
	}

	CheckErr(f.Close())
}

func GetPid() int {
	p, _ := ioutil.ReadFile("beans.pid")
	pid, _ := strconv.Atoi(string(p))
	return pid
}

func CheckErr(err error)  {
	if err != nil {
		panic(err)
	}
}

func Start()  {
	cmdName := "beans"

	if lp, err := exec.LookPath(cmdName); err != nil {
		fmt.Println("look path error:", err)
		os.Exit(1)
	} else {
		cmdName = lp
	}

	fmt.Println(cmdName)

	procAttr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	args := []string{cmdName}
	for key,value := range os.Args {
		if key == 0 {
			continue
		}
		args = append(args, value)
	}
	fmt.Println(args)

	process, err := os.StartProcess(cmdName, args, procAttr)

	RecordPid(process)
	if err != nil {
		fmt.Println("start process error:", err)
		os.Exit(2)
	}
	fmt.Println(process.Pid)
}