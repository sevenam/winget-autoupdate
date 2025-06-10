package main

import (
	"fmt"
	"log"
	"os"
)

func logMessage(message string) {
	fmt.Println(message)

	// Ensure the log directory exists
	if _, err := os.Stat(LogDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(LogDirectory, os.ModePerm)
		if err != nil {
			log.Fatal("Error creating log directory:", err)
		}
	}

	f, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(message)
}
