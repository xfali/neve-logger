// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package xlogneve

import (
	"fmt"
	"github.com/xfali/fig"
	"github.com/xfali/neve-core/bean"
	"github.com/xfali/stream"
	"github.com/xfali/xlog"
	"github.com/xfali/xlog/writer"
	"io"
	"os"
	"strings"
)

type processor struct {
	writers []io.WriteCloser

	formatter xlog.Formatter
	creator   FileWriterFactory
	level     xlog.Level
}

type logCaller struct {
	File string `json:"file" yaml:"file" `
	Func string `json:"func" yaml:"func" `
}
type logConf struct {
	Level        string    `json:"level" yaml:"level" `
	File         []string  `json:"file" yaml:"file" `
	LogCaller    logCaller `json:"caller" yaml:"caller" `
	SimpleName   bool      `json:"simpleName" yaml:"simpleName" `
	NoFatalTrace bool      `json:"noFatalTrace" yaml:"noFatalTrace" `
}

type Opt func(*processor)

type FileWriterFactory func(path string) (io.WriteCloser, error)

func NewLoggerProcessor(opts ...Opt) *processor {
	ret := &processor{
		creator: defaultCreator,
		level:   -1,
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}

func (p *processor) Init(conf fig.Properties, container bean.Container) error {
	var logConf logConf
	err := conf.GetValue("neve.logger", &logConf)
	if err != nil {
		xlog.Errorln("Get logger config failed.")
		return nil
	}

	lvStr := logConf.Level
	lv, _ := transLevel(lvStr)
	if p.level == -1 {
		p.level = lv
	}

	flag := parseCaller(logConf.LogCaller)
	logging := xlog.NewLogging(xlog.SetCallerFlag(flag),
		xlog.SetCallerFormatter(callerFormatter(logConf)),
		xlog.SetFatalNoTrace(logConf.NoFatalTrace))
	logging.SetSeverityLevel(p.level)

	if p.formatter != nil {
		logging.SetFormatter(p.formatter)
	}

	writers, err := p.parseWriter(logConf.File)
	if err != nil {
		// Init error, close all
		p.closeAll()
		xlog.Errorln(err)
		return err
	}
	if len(p.writers) > 0 {
		logging.SetOutput(io.MultiWriter(writers...))
	}

	fac := xlog.NewFactory(logging)
	if logConf.SimpleName {
		fac.SimplifyNameFunc = xlog.SimplifyNameFirstLetter
	}
	xlog.ResetFactory(fac)
	return nil
}

func callerFormatter(config logConf) func(file string, line int, funcName string) string {
	return func(file string, line int, funcName string) string {
		if funcName != "" && config.LogCaller.Func == "simple" {
			funcName = simpleFuncName(funcName)
		}
		if file != "" {
			if funcName == "" {
				return fmt.Sprintf("%s:%d", file, line)
			} else {
				return fmt.Sprintf("%s:%d (%s)", file, line, funcName)
			}
		} else {
			if funcName == "" {
				return ""
			} else {
				return "(" + funcName + ")"
			}
		}
	}
}

func simpleFuncName(funcName string) string {
	segs := strings.Split(funcName, "/")
	buf := strings.Builder{}
	buf.Grow(len(funcName) / 2)
	size := len(segs) - 1
	for i := 0; i < size; i++ {
		if len(segs[i]) > 0 {
			buf.WriteString(segs[i][:1])
			buf.WriteString(".")
		}
	}
	buf.WriteString(segs[size])
	return buf.String()
}

func parseCaller(caller logCaller) int {
	flag := xlog.CallerShortFile
	if caller.File == "none" {
		flag = xlog.CallerNone
	} else if caller.File == "long" {
		flag = xlog.CallerLongFile
	}

	if caller.Func == "short" {
		flag |= xlog.CallerShortFunc
	} else if caller.Func == "long" || caller.Func == "simple" {
		flag |= xlog.CallerLongFunc
	}
	return flag
}

func (p *processor) parseWriter(outputs []string) ([]io.Writer, error) {
	var writers []io.Writer
	errOp := stream.Slice(outputs).Distinct(func(s1, s2 string) bool {
		return s1 == s2
	}).Map(func(s string) error {
		w := matchOsOutput(s)
		if w != nil {
			writers = append(writers, w)
			return nil
		} else {
			w, err := p.creator(s)
			if err != nil {
				return err
			}
			writers = append(writers, w)
			// Add for close
			p.writers = append(p.writers, w)
			return nil
		}
	}).Filter(func(err error) bool {
		return err != nil
	}).FindFirst()

	if errOp.IsPresent() {
		return writers, errOp.Get().(error)
	} else {
		return writers, nil
	}
}

func matchOsOutput(op string) io.Writer {
	if len(op) == 0 {
		return nil
	}
	if strings.ToLower(op) == "stdout" {
		return os.Stdout
	} else if strings.ToLower(op) == "stderr" {
		return os.Stderr
	}
	return nil
}

func (p *processor) Classify(o interface{}) (bool, error) {
	return false, nil
}

func (p *processor) Process() error {
	return nil
}

func (p *processor) closeAll() error {
	var ret error
	if len(p.writers) > 0 {
		for _, w := range p.writers {
			err := w.Close()
			if err != nil {
				ret = err
			}
		}
		p.writers = nil
	}
	return ret
}

func (p *processor) BeanDestroy() error {
	return p.closeAll()
}

func defaultCreator(path string) (closer io.WriteCloser, e error) {
	w := writer.NewBufferedRotateFileWriter(&writer.BufferedRotateFile{
		Path:            path,
		RotateFrequency: writer.RotateEveryDay,
		RotateFunc:      writer.ZipLogsAsync,
	})
	if w == nil {
		return w, fmt.Errorf("Init logger failed, log file: %s. ", path)
	}
	return w, nil
}

func transLevel(s string) (xlog.Level, bool) {
	if len(s) > 0 {
		s = strings.ToUpper(s)
		switch s {
		case "DEBUG":
			return xlog.DEBUG, true
		case "INFO":
			return xlog.INFO, true
		case "WARN":
			return xlog.WARN, true
		case "ERROR":
			return xlog.ERROR, true
		case "PANIC":
			return xlog.PANIC, true
		case "FATAL":
			return xlog.FATAL, true
		}
	}
	return xlog.INFO, false
}

func OptSetFileWriterFactory(fac FileWriterFactory) Opt {
	return func(p *processor) {
		p.creator = fac
	}
}

func OptSetLogLevel(level xlog.Level) Opt {
	return func(p *processor) {
		p.level = level
	}
}

func OptSetLogFormatter(formatter xlog.Formatter) Opt {
	return func(p *processor) {
		p.formatter = formatter
	}
}
