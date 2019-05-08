#!/usr/bin/env bash
go install beanstalkf.go
go install beanstalkf-daemon.go
echo "The Command File Generate In $GOBIN"
ls -l -t $GOBIN