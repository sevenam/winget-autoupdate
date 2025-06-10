package main

import (
	"fmt"
	"log"
	"os"
)

func logMessage(message string) {
	fmt.Println(message)

	logDir := LogDirectory
	const logFile = LogFilePath
	// Ensure the log directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatal("Error creating log directory:", err)
		}
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(message)
}
