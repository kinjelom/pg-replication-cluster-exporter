package main

import (
	"fmt"
	"os"
	"time"
)

type Log struct {
	Verbosity int
}

func NewLog() *Log {
	return &Log{Verbosity: 1}
}

func (log *Log) info(f string, a ...interface{}) {
	if log.Verbosity >= 2 {
		_, _ = fmt.Fprintf(os.Stdout, "["+time.Now().Format(time.RFC3339)+"] INFO "+f+"\n", a...)
	}
}

func (log *Log) warn(f string, a ...interface{}) {
	if log.Verbosity >= 1 {
		_, _ = fmt.Fprintf(os.Stdout, "["+time.Now().Format(time.RFC3339)+"] WARNING "+f+"\n", a...)
	}
}

func (log *Log) error(f string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "["+time.Now().Format(time.RFC3339)+"] ERROR "+f+"\n", a...)
}

func (log *Log) debug(f string, a ...interface{}) {
	if log.Verbosity >= 3 {
		_, _ = fmt.Fprintf(os.Stderr, "["+time.Now().Format(time.RFC3339)+"] DEBUG "+f+"\n", a...)
	}
}

var log = NewLog()
