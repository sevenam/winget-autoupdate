package main

import (
	"fmt"
	"os/exec"
)

func installScheduledTask() {
	taskName := "WingetAutoUpdate"
	exePath := "D:\\git\\sevenam\\winget-autoupdate\\wingetau.exe"
	args := "update"

	cmd := exec.Command("schtasks",
		"/Create",
		"/F",
		"/SC", "MINUTE",
		"/MO", "1",
		"/TN", taskName,
		"/TR", fmt.Sprintf("\"%s %s\"", exePath, args),
		"/RL", "HIGHEST",
		"/RU", "SYSTEM",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to create scheduled task: %v\nOutput: %s\n", err, string(output))
		return
	}
	fmt.Println("Scheduled task created successfully.")
}

func uninstallScheduledTask() {
	taskName := "WingetAutoUpdate"
	cmd := exec.Command("schtasks", "/Delete", "/TN", taskName, "/F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to delete scheduled task: %v\nOutput: %s\n", err, string(output))
		return
	}
	fmt.Println("Scheduled task deleted successfully.")
}
