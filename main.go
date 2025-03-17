package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/gen2brain/beeep"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// the comment below is required!
//
//go:embed assets/*
var assets embed.FS

const serviceName = "Winget-AutoUpdate"
const logFile = "C:\\ProgramData\\winget-service\\winget-service.log"

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

func installService() {
	fmt.Println("Installing service...")

	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}

	m, err := mgr.Connect()
	if err != nil {
		log.Fatal("Failed to connect to service manager:", err)
	}
	defer m.Disconnect()

	s, err := m.CreateService(serviceName, exePath, mgr.Config{
		StartType: mgr.StartAutomatic,
	}, "runservice")
	if err != nil {
		log.Fatal("Failed to create service:", err)
	}
	defer s.Close()

	// Set the "Allow service to interact with desktop" option
	h, err := windows.OpenService(m.Handle, windows.StringToUTF16Ptr(serviceName), windows.SERVICE_CHANGE_CONFIG|windows.SERVICE_QUERY_CONFIG)
	if err != nil {
		log.Fatal("Failed to open service:", err)
	}
	defer windows.CloseServiceHandle(h)

	var bytesNeeded uint32
	err = windows.QueryServiceConfig(h, nil, 0, &bytesNeeded)
	if err != windows.ERROR_INSUFFICIENT_BUFFER {
		log.Fatal("Failed to query service config size:", err)
	}

	buffer := make([]byte, bytesNeeded)
	serviceConfig := (*windows.QUERY_SERVICE_CONFIG)(unsafe.Pointer(&buffer[0]))
	err = windows.QueryServiceConfig(h, serviceConfig, bytesNeeded, &bytesNeeded)
	if err != nil {
		log.Fatal("Failed to query service config:", err)
	}

	serviceType := serviceConfig.ServiceType | windows.SERVICE_INTERACTIVE_PROCESS
	err = windows.ChangeServiceConfig(h, serviceType, windows.SERVICE_NO_CHANGE, windows.SERVICE_NO_CHANGE, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		log.Fatal("Failed to change service config:", err)
	}

	fmt.Println("Service installed successfully.")

	startService()
}

func uninstallService() {
	fmt.Println("Uninstalling service...")
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal("Failed to connect to service manager:", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		log.Fatal("Service not found.")
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		log.Fatal("Failed to delete service:", err)
	}

	fmt.Println("Service uninstalled successfully.")
}

func startService() {
	fmt.Println("Starting service...")
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal("Failed to connect to service manager:", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		log.Fatal("Service not found:", err)
	}
	defer s.Close()

	err = s.Start("runservice")
	if err != nil {
		log.Fatal("Failed to start service:", err)
	}

	fmt.Println("Service started successfully.")
}

func stopService() {
	fmt.Println("Stopping service...")
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal("Failed to connect to service manager:", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		log.Fatal("Service not found.")
	}
	defer s.Close()

	status, err := s.Control(svc.Stop)
	if err != nil {
		log.Fatal("Failed to stop service:", err)
	}

	fmt.Printf("Service stopping (status: %v)...\n", status.State)
}

func logMessage(message string) {
	fmt.Println(message)

	// Ensure the log directory exists
	logDir := "C:\\ProgramData\\winget-service"
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

func runAsService() {
	fmt.Println("Running as a Windows Service...")
	run := svc.Run
	if err := run(serviceName, &wingetauWindowsService{}); err != nil {
		log.Fatalf("Failed to run service: %v", err)
	}
}
