#!/usr/bin/env bash
go get -t github.com/chenbo29/go-beanstalkd-client
go build beanstalkf-order.go
go build beanstalkf-daemon.go
cp beanstalkf-order /usr/bin/beanstalkd-order
cp beanstalkf-daemon /usr/bin/beanstalkd-daemon
echo "beanstalkf-order and beanstalkf-daemon install success"
