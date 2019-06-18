package main

import (
	"flag"
	"path/filepath"
)

var sourceVolPath string
var logPath string
var logFilename string

func init() {
	flag.StringVar(&sourceVolPath, "source-path", "/tmp/angmar/src", "`location` where all source repositories are located")
	flag.StringVar(&logPath, "log-path", "/tmp/angmar/log", "`location` where all source repositories are located")
	flag.StringVar(&logFilename, "log-filename", "angmar.log", "`filename` for logs")
}

func getLogfileName() string {
	return filepath.Join(logPath, logFilename)
}
