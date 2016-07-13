// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"fmt"
	"io"
	"log"

	"github.com/admpub/core"
)

type CLogger struct {
	Name      string
	Disabled  bool
	Processor func(tag string, format string, args []interface{}) (string, []interface{})
	Logger    core.ILogger
}

func (c *CLogger) Error(v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		_, v = c.Processor(c.Name, ``, v)
		if v == nil {
			return
		}
	}
	c.Logger.Error(v...)
}

func (c *CLogger) Errorf(format string, v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		format, v = c.Processor(c.Name, format, v)
		if v == nil {
			return
		}
	}
	c.Logger.Errorf(format, v...)
}

func (c *CLogger) Debug(v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		_, v = c.Processor(c.Name, ``, v)
		if v == nil {
			return
		}
	}
	c.Logger.Debug(v...)
}

func (c *CLogger) Debugf(format string, v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		format, v = c.Processor(c.Name, format, v)
		if v == nil {
			return
		}
	}
	c.Logger.Debugf(format, v...)
}

func (c *CLogger) Info(v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		_, v = c.Processor(c.Name, ``, v)
		if v == nil {
			return
		}
	}
	c.Logger.Info(v...)
}

func (c *CLogger) Infof(format string, v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		format, v = c.Processor(c.Name, format, v)
		if v == nil {
			return
		}
	}
	c.Logger.Infof(format, v...)
}

func (c *CLogger) Warn(v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		_, v = c.Processor(c.Name, ``, v)
		if v == nil {
			return
		}
	}
	c.Logger.Warn(v...)
}

func (c *CLogger) Warnf(format string, v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		format, v = c.Processor(c.Name, format, v)
		if v == nil {
			return
		}
	}
	c.Logger.Warnf(format, v...)
}

type TLogger struct {
	SQL   *CLogger
	Event *CLogger
	Cache *CLogger
	ETime *CLogger
	Base  *CLogger
	Other *CLogger
}

func (t *TLogger) Open(tags ...string) {
	if len(tags) == 0 {
		t.SQL.Disabled = false
		t.Event.Disabled = false
		t.Cache.Disabled = false
		t.ETime.Disabled = false
		t.Base.Disabled = false
		t.Other.Disabled = false
		return
	}
	for _, tag := range tags {
		t.SetStatusByName(tag, true)
	}
}

func (t *TLogger) Close(tags ...string) {
	if len(tags) == 0 {
		t.SQL.Disabled = true
		t.Event.Disabled = true
		t.Cache.Disabled = true
		t.ETime.Disabled = true
		t.Base.Disabled = true
		t.Other.Disabled = true
		return
	}
	for _, tag := range tags {
		t.SetStatusByName(tag, false)
	}
}

func (t *TLogger) SetStatusByName(tag string, status bool) {
	switch tag {
	case "sql":
		t.SQL.Disabled = status
	case "event":
		t.Event.Disabled = status
	case "cache":
		t.Cache.Disabled = status
	case "etime":
		t.ETime.Disabled = status
	case "base":
		t.Base.Disabled = status
	case "other":
		t.Other.Disabled = status
	}
}

func (t *TLogger) SetLogger(logger core.ILogger) {
	t.SQL.Logger = logger
	t.Event.Logger = logger
	t.Cache.Logger = logger
	t.ETime.Logger = logger
	t.Base.Logger = logger
	t.Other.Logger = logger
}

func NewTLogger(logger core.ILogger) *TLogger {
	return &TLogger{
		SQL:   &CLogger{Name: "sql", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		Event: &CLogger{Name: "event", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		Cache: &CLogger{Name: "cache", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		ETime: &CLogger{Name: "etime", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		Base:  &CLogger{Name: "base", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		Other: &CLogger{Name: "other", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
	}
}

var defaultLogProcessor = func(tag string, format string, args []interface{}) (string, []interface{}) {
	if format == "" {
		if len(args) > 0 {
			args[0] = fmt.Sprintf("[%s] %v", tag, args[0])
		}
		return format, args
	}
	format = "[" + tag + "] " + format
	return format, args
}

const (
	DEFAULT_LOG_PREFIX = "[dbx]"
	DEFAULT_LOG_FLAG   = log.Ldate | log.Lmicroseconds
	DEFAULT_LOG_LEVEL  = core.LOG_DEBUG
)

var _ core.ILogger = DiscardLogger{}

type DiscardLogger struct{}

func (DiscardLogger) Debug(v ...interface{})                 {}
func (DiscardLogger) Debugf(format string, v ...interface{}) {}
func (DiscardLogger) Error(v ...interface{})                 {}
func (DiscardLogger) Errorf(format string, v ...interface{}) {}
func (DiscardLogger) Info(v ...interface{})                  {}
func (DiscardLogger) Infof(format string, v ...interface{})  {}
func (DiscardLogger) Warn(v ...interface{})                  {}
func (DiscardLogger) Warnf(format string, v ...interface{})  {}
func (DiscardLogger) Level() core.LogLevel {
	return core.LOG_UNKNOWN
}
func (DiscardLogger) SetLevel(l core.LogLevel) {}
func (DiscardLogger) ShowSQL(show ...bool)     {}
func (DiscardLogger) IsShowSQL() bool {
	return false
}

// SimpleLogger is the default implment of core.ILogger
type SimpleLogger struct {
	DEBUG   *log.Logger
	ERR     *log.Logger
	INFO    *log.Logger
	WARN    *log.Logger
	level   core.LogLevel
	showSQL bool
}

var _ core.ILogger = &SimpleLogger{}

// NewSimpleLogger use a special io.Writer as logger output
func NewSimpleLogger(out io.Writer) *SimpleLogger {
	return NewSimpleLogger2(out, DEFAULT_LOG_PREFIX, DEFAULT_LOG_FLAG)
}

// NewSimpleLogger2 let you customrize your logger prefix and flag
func NewSimpleLogger2(out io.Writer, prefix string, flag int) *SimpleLogger {
	return NewSimpleLogger3(out, prefix, flag, DEFAULT_LOG_LEVEL)
}

// NewSimpleLogger3 let you customrize your logger prefix and flag and logLevel
func NewSimpleLogger3(out io.Writer, prefix string, flag int, l core.LogLevel) *SimpleLogger {
	return &SimpleLogger{
		DEBUG: log.New(out, fmt.Sprintf("%s [debug] ", prefix), flag),
		ERR:   log.New(out, fmt.Sprintf("%s [error] ", prefix), flag),
		INFO:  log.New(out, fmt.Sprintf("%s [info]  ", prefix), flag),
		WARN:  log.New(out, fmt.Sprintf("%s [warn]  ", prefix), flag),
		level: l,
	}
}

// Error implement core.ILogger
func (s *SimpleLogger) Error(v ...interface{}) {
	if s.level <= core.LOG_ERR {
		s.ERR.Output(2, fmt.Sprint(v...))
	}
	return
}

// Errorf implement core.ILogger
func (s *SimpleLogger) Errorf(format string, v ...interface{}) {
	if s.level <= core.LOG_ERR {
		s.ERR.Output(2, fmt.Sprintf(format, v...))
	}
	return
}

// Debug implement core.ILogger
func (s *SimpleLogger) Debug(v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		s.DEBUG.Output(2, fmt.Sprint(v...))
	}
	return
}

// Debugf implement core.ILogger
func (s *SimpleLogger) Debugf(format string, v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		s.DEBUG.Output(2, fmt.Sprintf(format, v...))
	}
	return
}

// Info implement core.ILogger
func (s *SimpleLogger) Info(v ...interface{}) {
	if s.level <= core.LOG_INFO {
		s.INFO.Output(2, fmt.Sprint(v...))
	}
	return
}

// Infof implement core.ILogger
func (s *SimpleLogger) Infof(format string, v ...interface{}) {
	if s.level <= core.LOG_INFO {
		s.INFO.Output(2, fmt.Sprintf(format, v...))
	}
	return
}

// Warn implement core.ILogger
func (s *SimpleLogger) Warn(v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		s.WARN.Output(2, fmt.Sprint(v...))
	}
	return
}

// Warnf implement core.ILogger
func (s *SimpleLogger) Warnf(format string, v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		s.WARN.Output(2, fmt.Sprintf(format, v...))
	}
	return
}

// Level implement core.ILogger
func (s *SimpleLogger) Level() core.LogLevel {
	return s.level
}

// SetLevel implement core.ILogger
func (s *SimpleLogger) SetLevel(l core.LogLevel) {
	s.level = l
	return
}
