package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	lockLimit     = flag.Int("lock-limit", 5, "Lock screen after this many skips/postpones")
	shutdownLimit = flag.Int("shutdown-limit", 10, "Shutdown after this many skips/postpones")
	curfewStart   = flag.Int("curfew-start", 22, "Curfew start hour, 24h local time")
	curfewEnd     = flag.Int("curfew-end", 8, "Curfew end hour, 24h local time")
	curfewRelock  = flag.Duration("curfew-relock", time.Minute, "Re-lock interval during curfew")
)

func main() {
	flag.Parse()

	printBanner(*lockLimit, *shutdownLimit)

	events, err := StartMonitor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start monitor: %v\n", err)
		os.Exit(1)
	}

	enforcer := NewEnforcer(*lockLimit, *shutdownLimit)

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to resolve home dir: %v\n", err)
		os.Exit(1)
	}
	go StartCurfew(*curfewStart, *curfewEnd, *curfewRelock, filepath.Join(home, ".defixate-curfew-off"))

	logMsg("Monitoring Time Out logs...")
	logMsg("Press Ctrl+C to stop")
	fmt.Println()

	for event := range events {
		switch event {
		case EventPostponed:
			logMsg("Detected: Postpone")
			enforcer.OnSkipOrPostpone()
		case EventSkipped:
			logMsg("Detected: Skip")
			enforcer.OnSkipOrPostpone()
		case EventCompleted:
			logMsg("Detected: Break finished")
			enforcer.OnBreakCompleted()
		}
	}
}

func printBanner(lockLimit, shutdownLimit int) {
	fmt.Println()
	fmt.Printf("\033[1;36m╔═══════════════════════════════════════════════════════════╗\033[0m\n")
	fmt.Printf("\033[1;36m║  TIMEOUT-DEFIXATE - %d skips = lock, %d = shutdown         ║\033[0m\n", lockLimit, shutdownLimit)
	fmt.Printf("\033[1;36m║  Using Time Out's explicit log messages                   ║\033[0m\n")
	fmt.Printf("\033[1;36m╚═══════════════════════════════════════════════════════════╝\033[0m\n")
	fmt.Println()
	os.Stdout.Sync()
}

func logMsg(msg string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("\033[33m%s │ %s\033[0m\n", timestamp, msg)
	os.Stdout.Sync()
}

func logSuccess(msg string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("\033[32m%s │ %s\033[0m\n", timestamp, msg)
	os.Stdout.Sync()
}

func logWarn(msg string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("\033[31m%s │ %s\033[0m\n", timestamp, msg)
	os.Stdout.Sync()
}
