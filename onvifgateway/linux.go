// +build linux

package main

import (
	"syscall"

	"wholeally.com/common/logs"
)

func init() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if nil != err {
		logs.Error(err)
	}

	rLimit.Cur = 1000000
	rLimit.Max = 1000000
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if nil != err {
		logs.Error(err)
	} else {
		logs.Info("setrlimit: ", rLimit.Max)
	}
}
