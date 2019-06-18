package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/step/angmar/pkg/tarutils"
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

func getLogfile() *os.File {
	file, err := os.OpenFile(getLogfileName(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return file
}

func getExtractorGenerator() tarutils.ExtractorGenerator {
	return tarutils.DefaultExtractorGenerator{Src: sourceVolPath}
}
