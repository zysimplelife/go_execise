# go_execise

## How to build example
* Download code
```
go get github.com/zysimplelife/go_execise
```
* Download vendor
```
cd $GOPATH/src/github.com/zysimplelife/go_execise/DIRECTORY_YOU_WANT
dep init
```
* Build
```
go build
```

## secert
Experiment on how to use go to monitor an certifcate and reload service.
It produces new certifcate every 10 second, which will trigger the restart of service

## pipelines
Experiment on how to stop multipl go routing with the help of channel and WaitGroup
It is a common pattern descriped in [Go Concurrency Patterns: Pipelines and cancellation](https://blog.golang.org/pipelines).
