package main

import (
	"flag"

	"go.uber.org/zap"
)

// TODO: Think about config file for logger too
func main() {
	logLevel = zap.LevelFlag("loglevel", zap.InfoLevel, "set logging level")
	var filelog string
	flag.StringVar(&filelog, "filelog", "", "choose file for logs, leave empty to use stderr")
	flag.Parse()
}
