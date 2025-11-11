package main

import (
	"os"

	"github.com/charmbracelet/log"
)

func getNewLogger(prefix string) *log.Logger {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return log.NewWithOptions(file, log.Options{
		TimeFormat:      "01:04:05.000",
		Level:           log.DebugLevel,
		ReportTimestamp: true,
		Prefix:          prefix,
		// ReportCaller:    true,
	})
}
