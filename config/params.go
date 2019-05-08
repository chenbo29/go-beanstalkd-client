package config

import (
	"flag"
	"fmt"
	"os"
)

// 应用的执行参数配置信息
type ParamsData struct {
	Name        string
	Description string
	Host        string
	Port        string
	Tube        string
	IsAllStatus bool
	Daemon      bool
}

var returnParams ParamsData

// init 初始化命令执行的参数说明和基础配置
func init() {
	returnParams.Name = "go-beanstalk-client"
	returnParams.Description = "go-beanstalk-client By chenbotome@163.com"
	returnParams.Host = "127.0.0.1"
	returnParams.Port = "11300"
	if len(os.Args) > 1 {
		CommandLine := flag.NewFlagSet(os.Args[1], 0)
		CommandLine.Usage = func() {
			fmt.Fprintf(os.Stderr, "go-beanstalkd created by chenbotome@163.com \n Usage of %s:\n", os.Args[0])
			CommandLine.PrintDefaults()
			os.Exit(0)
		}
		switch os.Args[1] {
		case "start":
			CommandLine.StringVar(&returnParams.Port, "port", "11300", "the port of beanstalkd")
			CommandLine.StringVar(&returnParams.Port, "p", "11300", "the port of beanstalkd (shorthand)")
			CommandLine.StringVar(&returnParams.Host, "host", "127.0.0.1", "the host of beanstalkd")
			CommandLine.StringVar(&returnParams.Host, "h", "127.0.0.1", "the host of beanstalkd (shorthand)")
			CommandLine.BoolVar(&returnParams.Daemon, "daemon", false, "Start With Daemon")
			CommandLine.BoolVar(&returnParams.Daemon, "d", false, "Start With Daemon (shorthand)")
			CommandLine.Parse(os.Args[2:])
		case "status":
			CommandLine.StringVar(&returnParams.Tube, "tube", "default", "the status of tube")
			CommandLine.StringVar(&returnParams.Tube, "t", "default", "the status of tube (shorthand)")
			CommandLine.BoolVar(&returnParams.IsAllStatus, "all", false, "the status of beanstalk")
			CommandLine.BoolVar(&returnParams.IsAllStatus, "a", false, "the status of beanstalk (shorthand)")
			CommandLine.Parse(os.Args[2:])
		}
	}
}

// GetParams 获取配置参数信息
func GetParams() *ParamsData {
	return &returnParams
}
