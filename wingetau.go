package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const serviceName = "WinGetUpdater"
const logFile = "C:\\ProgramData\\winget-service\\winget-service.log"

type WingetPackage struct {
	Id               string `json:"PackageIdentifier"`
	Name             string `json:"PackageName"`
	Version          string `json:"Version"`
	AvailableVersion string `json:"AvailableVersion"`
	Source           string `json:"Source"`
}

func main() {
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

	fmt.Println("Service installed successfully.")
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
		log.Fatal("Service not found.")
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

func runWinGetUpdate() {
	fmt.Println("Checking for package updates...")

	cmd := exec.Command("winget", "list", "--upgrade-available", "--source=winget")
	output, err := cmd.Output()
	if err != nil {
		logMessage(fmt.Sprintf("Error checking updates: %v", err))
		sendNotification("WinGet Update", fmt.Sprintf("Error checking updates: %v", err), "error")
		return
	}

	packages := parseWingetOutput(string(output))

	if len(packages) == 0 {
		logMessage("No updates available.")
		sendNotification("WinGet Update", "No updates available.", "info")
		return
	}

	for _, pkg := range packages {
		logMessage(fmt.Sprintf("Updating %s (%s -> %s)...", pkg.Name, pkg.Version, pkg.AvailableVersion))
		sendNotification("WinGet Update", fmt.Sprintf("Updating %s (%s -> %s)...", pkg.Name, pkg.Version, pkg.AvailableVersion), "info")
		updateCmd := exec.Command("winget", "upgrade", "--silent", "--id", pkg.Id)
		updateOutput, updateErr := updateCmd.CombinedOutput()

		if updateErr != nil {
			logMessage(fmt.Sprintf("Error updating %s: %v\nOutput: %s", pkg.Name, updateErr, string(updateOutput)))
			sendNotification("WinGet Update", fmt.Sprintf("Error updating %s: %v\nOutput: %s", pkg.Name, updateErr, string(updateOutput)), "error")
		} else {
			logMessage(fmt.Sprintf("Successfully updated %s to version %s", pkg.Name, pkg.AvailableVersion))
			sendNotification("WinGet Update", fmt.Sprintf("Successfully updated %s to version %s", pkg.Name, pkg.AvailableVersion), "success")
		}
	}
}

func parseWingetOutput(output string) []WingetPackage {
	var packages []WingetPackage
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) == 0 || strings.Count(trimmedLine, "-") == len(trimmedLine) {
			continue
		}
		if strings.Contains(line, "Name") && strings.Contains(line, "Id") && strings.Contains(line, "Version") && strings.Contains(line, "Available") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		pkg := WingetPackage{
			Name:             strings.Join(fields[:len(fields)-3], " "),
			Id:               fields[len(fields)-3],
			Version:          fields[len(fields)-2],
			AvailableVersion: fields[len(fields)-1],
		}
		packages = append(packages, pkg)
	}
	return packages
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
		iconPath = ".\\error.png" // Replace with the path to your error icon
	case "success":
		iconPath = ".\\success.png" // Replace with the path to your success icon
	case "info":
		iconPath = ".\\info.png" // Replace with the path to your info icon
	default:
		iconPath = ""
	}

	err := beeep.Notify(title, message, iconPath)
	if err != nil {
		logMessage(fmt.Sprintf("Failed to send notification: %v", err))
	}
}

func runAsService() {
	fmt.Println("Running as a Windows Service...")
	for {
		runWinGetUpdate()
		time.Sleep(24 * time.Hour)
	}
}
