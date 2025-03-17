package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type wingetService struct{}

func runWinGetUpdate() {
	fmt.Println("Checking for package updates...")

	wingetPath, err := findWingetPath()
	if err != nil {
		logMessage(fmt.Sprintf("Error finding winget path: %v", err))
		sendNotification("WinGet Update", fmt.Sprintf("Error finding winget path: %v", err), "error")
		return
	}

	updateAgreements(wingetPath)

	cmd := exec.Command(wingetPath, "list", "--upgrade-available", "--source=winget")
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
		updateCmd := exec.Command(wingetPath, "upgrade", "--silent", "--include-unknown", "--accept-package-agreements", "--accept-source-agreements", "--disable-interactivity", "--id", pkg.Id)
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

func updateAgreements(wingetPath string) {
	// update msstore and agree to the package agreements
	cmd := exec.Command(wingetPath, "source", "update", "--name", "msstore", "--disable-interactivity")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logMessage(fmt.Sprintf("Error updating msstore source: %v", err))
		logMessage(fmt.Sprintf("Output: %s", string(output)))
		sendNotification("WinGet Update", fmt.Sprintf("Error updating msstore source: %v", err), "error")
		return
	}
	logMessage(fmt.Sprintf("Output: %s", string(output)))
}

func findWingetPath() (string, error) {
	basePath := "C:\\Program Files\\WindowsApps"
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, "winget.exe") {
			return fmt.Errorf(path) // Use error to return the path
		}
		return nil
	})
	if err != nil {
		if strings.HasSuffix(err.Error(), "winget.exe") {
			return err.Error(), nil // Extract the path from the error
		}
		return "", err
	}
	return "", fmt.Errorf("winget.exe not found")
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
