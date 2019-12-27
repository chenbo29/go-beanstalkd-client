#!/usr/bin/env bash
go build beanstalkf-post.go
go build beanstalkf-post-daemon.go
cp beanstalkf-post /usr/bin/beanstalkd-post
cp beanstalkf-post-daemon /usr/bin/beanstalkd-post-daemon
echo "beanstalkf-post and beanstalkf-post-daemon install success"
