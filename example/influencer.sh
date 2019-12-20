#!/usr/bin/env bash
go get -t github.com/chenbo29/go-beanstalkd-client
go build beanstalkf-influencer.go
go build beanstalkf-daemon.go
cp beanstalkf-influencer /usr/bin/beanstalkd-influencer
cp beanstalkf-daemon /usr/bin/beanstalkd-daemon
echo "beanstalkf-influencer and beanstalkf-daemon install success"
