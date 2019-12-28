#!/usr/bin/env bash
go get -u github.com/chenbo29/go-beanstalkd-client
go build beanstalkf-influencer.go
go build beanstalkf-influencer-daemon.go
cp beanstalkf-influencer /usr/bin/beanstalkf-influencer
cp beanstalkf-influencer-daemon /usr/bin/beanstalkf-influencer-daemon
echo "beanstalkf-influencer and beanstalkf-influencer-daemon install success"
