// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/neve-core"
	"github.com/xfali/neve-logger/xlogneve"
	"github.com/xfali/xlog"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	app := neve.NewFileConfigApplication("assets/config-test.yaml")
	//app.RegisterBean(xlogneve.NewLoggerProcessor(xlogneve.OptSetLogFormatter(&xlog.JsonFormatter{})))
	app.RegisterBean(xlogneve.NewLoggerProcessor())
	go app.Run()

	time.Sleep(1 * time.Second)

	xlog.Infoln("cannot output, display this means failed")
	xlog.Warnln("this is ok")

	time.Sleep(2 * time.Second)
}

