package main

import (
	"fmt"
	"os"

	"github.com/gen2brain/beeep"
)

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
