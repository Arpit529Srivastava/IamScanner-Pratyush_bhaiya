package main

import "os"

type CustomLogger struct {
	Output *os.File
}

func (l *CustomLogger) Write(p []byte) (n int, err error) {
	return l.Output.Write(p)
}
