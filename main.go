package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/gen2brain/beeep"
	"golang.org/x/sys/windows/svc"
)

// the comment below is required - don't delete!
//
//go:embed assets/*
var assets embed.FS

const serviceName = ServiceName
const logFile = LogFilePath

func main() {
	isWindowsService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("Failed to determine if we are running as Windows Service: %v", err)
	}

	if isWindowsService {
		logMessage("Running as a Windows Service...")
		runAsService()
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: mycli <install|uninstall|start|stop|update|runservice>")
		return
	}

	switch os.Args[1] {
	case "install":
		installService()
	case "uninstall":
		uninstallService()
	case "start":
		startService()
	case "stop":
		stopService()
	case "update":
		runWinGetUpdate()
	case "runservice":
		runAsService()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func logMessage(message string) {
	fmt.Println(message)

	// Ensure the log directory exists
	logDir := LogDirectory
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

func sendNotification(title, message, notificationType string) {
	var iconPath string
	switch notificationType {
	case "error":
		iconPath = "assets/error.png" // Path to embedded error icon
	case "success":
		iconPath = "assets/success.png" // Path to embedded success icon
	case "info":
		iconPath = "assets/info.png" // Path to embedded info icon
	default:
		iconPath = ""
	}

	iconData, err := assets.ReadFile(iconPath)
	if err != nil {
		logMessage(fmt.Sprintf("Failed to read icon file: %v", err))
		iconPath = ""
	} else {
		// Write the icon data to a temporary file
		tmpFile, err := os.CreateTemp("", "icon-*.png")
		if err != nil {
			logMessage(fmt.Sprintf("Failed to create temp file for icon: %v", err))
			iconPath = ""
		} else {
			defer os.Remove(tmpFile.Name())
			_, err = tmpFile.Write(iconData)
			if err != nil {
				logMessage(fmt.Sprintf("Failed to write icon data to temp file: %v", err))
				iconPath = ""
			} else {
				iconPath = tmpFile.Name()
			}
		}
	}

	// Use the beeep package to send a notification to the Windows notification bar
	err = beeep.Notify(title, message, iconPath)
	if err != nil {
		logMessage(fmt.Sprintf("Failed to send notification: %v", err))
	}
}
