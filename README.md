# go-beanstalk-client
* [中文]("/README.md")
* [English]("/README.md")
## 目的
* 学习阅读GoLand的相关语法知识之后进行的练手项目，帮助自身理解基础语法并深入理解其语言的特性。
* 使用beanstalk队列（Beanstalk is a simple, fast work queue. https://beanstalkd.github.io/）并利用Go的协程goroutine特性，可以对每一个Tube管道创建一个工厂，工厂拥有若干自定义数量的工人，工人领取任务处理。
## 使用
* go run beans.go status -h
* go run beans.go start -h
## TODO