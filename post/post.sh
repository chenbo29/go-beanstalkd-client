#!/usr/bin/env bash
go build beanstalkf-post.go
go build beanstalkf-post-daemon.go
cp beanstalkf-post /usr/bin/beanstalkf-post
cp beanstalkf-post-daemon /usr/bin/beanstalkf-post-daemon
echo "beanstalkf-post and beanstalkf-post-daemon install success"
