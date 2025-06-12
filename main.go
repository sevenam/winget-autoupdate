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
		printHelp()
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
	case "help":
		printHelp()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func printHelp() {
	fmt.Printf(
		`wingetau.exe (Winget Auto Update) - A CLI with Windows Service capabilities for automatic updates of installed packages using WinGet.

Usage:
  %s <install|uninstall|start|stop|update|list|runservice>

Available commands:
  install      Install the service
  uninstall    Uninstall the service
  start        Start the service
  stop         Stop the service
  update       Run updates now
  list         List packages that would be updated
  runservice   Run as a Windows service
  help         Show this help message

`, filepath.Base(os.Args[0]))
}
