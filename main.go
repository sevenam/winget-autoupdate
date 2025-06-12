package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc"
)

// the comment below is required - don't delete!
//
//go:embed assets/*
var assets embed.FS

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
		fmt.Printf(
			`wingetau (Winget Auto Update) - A CLI with Windows Service capabilities for automatic updates of installed packages using WinGet.

Usage: %s <install|uninstall|start|stop|update|list|runservice>

Commands:
install    - Install the Windows service
uninstall  - Uninstall the Windows service
start      - Start the Windows service
stop       - Stop the Windows service
update     - Check for and apply updates to installed packages
list       - List available updates for installed packages
runservice - Run the application as a Windows Service

`,
			filepath.Base(os.Args[0]))
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
	case "list":
		listUpdates()
	case "runservice":
		runAsService()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
