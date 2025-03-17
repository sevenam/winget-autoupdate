package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/sys/windows/svc"
)

type wingetauWindowsService struct{}

func (m *wingetauWindowsService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	s <- svc.Status{State: svc.StartPending}

	// Report running status
	s <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	// Get the update interval from the environment variable
	intervalStr := os.Getenv("WINGET_UPDATE_INTERVAL_SECONDS")
	if intervalStr == "" {
		intervalStr = "60" // Default to 1 minute if not set
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		logMessage(fmt.Sprintf("Invalid update interval: %v", err))
		interval = 60 // Default to 1 minute if invalid
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			runWinGetUpdate()
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				break loop
			default:
				logMessage(fmt.Sprintf("Unexpected control request #%d", c))
			}
		}
	}

	s <- svc.Status{State: svc.StopPending}
	return false, 0
}
