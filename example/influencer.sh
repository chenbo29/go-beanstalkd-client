#!/usr/bin/env bash
export GOBIN=/usr/bin
go get -t github.com/chenbo29/go-beanstalkd-client
go build beanstalkf-order.go
go build beanstalkf-daemon.go
echo "beanstalkf-order and beanstalkf-daemon install success"
