package main

import (
	"fmt"
	"loggerSystemLLD/logger"
	"os"
)

func main() {
	file, err := os.OpenFile(
		"app.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open log file: %v\n", err)
		return
	}
	defer file.Close()

	log, err := logger.New(
		logger.InfoLevel,
		os.Stdout,
		file,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create logger: %v\n", err)
		return
	}

	if err := log.Debug("debug details"); err != nil {
		fmt.Fprintf(os.Stderr, "log debug: %v\n", err)
	}

	if err := log.Info("application started"); err != nil {
		fmt.Fprintf(os.Stderr, "log info: %v\n", err)
	}

	if err := log.Error("payment failed"); err != nil {
		fmt.Fprintf(os.Stderr, "log error: %v\n", err)
	}
}
