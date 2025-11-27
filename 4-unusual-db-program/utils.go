package main

import (
	"os"

	"github.com/charmbracelet/log"
)

func getNewLogger(prefix string, level log.Level) *log.Logger {
	return log.NewWithOptions(os.Stderr, log.Options{
		TimeFormat:      "01:04:05.000",
		Level:           level,
		ReportTimestamp: true,
		Prefix:          prefix,
	})
}
