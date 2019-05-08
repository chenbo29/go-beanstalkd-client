#!/usr/bin/env bash
go get github.com/chenbo29/go-beanstalkd-client
gopath=`go env GOPATH`
cd $gopath/src/github.com/chenbo29/go-beanstalkd-client/example
go install beanstalkf.go
go install beanstalkf-daemon.go
sudo cp $gopath/beanstalkf /bin/beanstalkf
sudo cp $gopath/beanstalkf-daemon /bin/beanstalkf-daemon
echo "beanstalkf and beanstalkf-daemon install success"
