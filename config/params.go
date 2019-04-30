package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var host string
var port string

type ParamsData struct {
	Host        string
	Port        string
	Num         uint64
	Tube        string
	IsAllStatus bool
	Daemon      bool
}

var returnParams ParamsData

var CommandLine *flag.FlagSet
var ErrHelp = errors.New("flag: help requested")

func init() {
	fmt.Println("init")
	returnParams.Host = "127.0.0.1"
	returnParams.Port = "11300"
	if len(os.Args) > 1 {
		CommandLine = flag.NewFlagSet(os.Args[1], 0)
		CommandLine.Usage = func() {
			fmt.Fprintf(os.Stderr, "go-beanstalkd created by chenbotome@163.com \n Usage of %s:\n", os.Args[0])
			CommandLine.PrintDefaults()
		}
		switch os.Args[1] {
		case "start":
			CommandLine.StringVar(&returnParams.Port, "port", "11300", "the port of beanstalkd")
			CommandLine.StringVar(&returnParams.Port, "p", "11300", "the port of beanstalkd (shorthand)")
			CommandLine.StringVar(&returnParams.Host, "host", "127.0.0.1", "the host of beanstalkd")
			CommandLine.StringVar(&returnParams.Host, "h", "127.0.0.1", "the host of beanstalkd (shorthand)")
			CommandLine.Uint64Var(&returnParams.Num, "num", 2, "the host of beanstalkd (shorthand)")
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

func GetParams() *ParamsData {
	return &returnParams
}
