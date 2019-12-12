#!/usr/bin/env bash
go install beanstalkf-influencer.go
go install beanstalkf-daemon.go
sudo cp $GOPATH/beanstalkf-influencer /usr/bin/beanstalkf-influencer
sudo cp $GOPATH/beanstalkf-daemon /usr/bin/beanstalkf-daemon
echo "beanstalkf-influencer and beanstalkf-daemon install success"
