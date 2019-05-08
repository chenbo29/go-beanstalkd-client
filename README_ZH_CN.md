# go-beanstalk-client
[![Go Report Card](https://goreportcard.com/badge/github.com/chenbo29/go-beanstalkd-client)](https://goreportcard.com/report/github.com/chenbo29/go-beanstalkd-client)
* [中文](/README_ZH_CN.md)
* [English](/README.md)
## 目的
* 学习阅读GoLand的相关语法知识之后进行的练手项目，帮助自身理解基础语法并深入理解其语言的特性。
* 使用beanstalk队列（Beanstalk is a simple, fast work queue. https://beanstalkd.github.io/ ）并利用Go的协程goroutine特性，可以对每一个Tube管道创建一个工厂，工厂拥有若干自定义数量的工人，工人领取任务处理。
## 安装-Linux
切换到example目录，执行安装脚本
```bash
./install.sh
```
## 安装-Windows
```bash
go install beanstalkf.go
go install beanstalkf-daemon.go
```
## 可执行文件
* beanstalkf
* beanstalkf-daemon 后台启动
## Command
* start
* status
* stop