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
		fmt.Printf("Usage: %s <install|uninstall|start|stop|update|runservice>\n", filepath.Base(os.Args[0]))
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
