package main

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

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

func runAsService() {
	fmt.Println("Running as a Windows Service...")
	run := svc.Run
	if err := run(serviceName, &windowsServiceRunner{}); err != nil {
		log.Fatalf("Failed to run service: %v", err)
	}
}
