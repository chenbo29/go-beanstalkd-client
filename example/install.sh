#!/usr/bin/env bash
go install beanstalkf.go
go install beanstalkf-daemon.go
sudo cp $GOBIN/beanstalkf /bin/beanstalkf
sudo cp $GOBIN/beanstalkf-daemon /bin/beanstalkf-daemon
echo "beanstalkf and beanstalkf-daemon install success"
