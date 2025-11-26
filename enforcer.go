package main

import (
	"fmt"
	"time"
)

type Enforcer struct {
	skipCount     int
	lockLimit     int
	shutdownLimit int
}

func NewEnforcer(lockLimit, shutdownLimit int) *Enforcer {
	return &Enforcer{
		skipCount:     0,
		lockLimit:     lockLimit,
		shutdownLimit: shutdownLimit,
	}
}

func (e *Enforcer) OnSkipOrPostpone() {
	e.skipCount++

	if e.skipCount >= e.shutdownLimit {
		logWarn(fmt.Sprintf("FINAL LIMIT! (%d/%d) - SHUTTING DOWN", e.skipCount, e.shutdownLimit))
		e.doShutdown()
	} else if e.skipCount >= e.lockLimit {
		untilShutdown := e.shutdownLimit - e.skipCount
		logWarn(fmt.Sprintf("LOCK! (%d/%d) - %d until shutdown", e.skipCount, e.shutdownLimit, untilShutdown))
		e.doLockLoop()
	} else {
		untilLock := e.lockLimit - e.skipCount
		logWarn(fmt.Sprintf("Skip/Postpone (%d/%d) - %d until lock", e.skipCount, e.shutdownLimit, untilLock))
	}
}

func (e *Enforcer) OnBreakCompleted() {
	logSuccess("Break completed! Resetting count.")
	e.skipCount = 0
}

func (e *Enforcer) doLockLoop() {
	fmt.Println()
	fmt.Printf("\033[1;31m╔══════════════════════════════════════════════════════════╗\033[0m\n")
	fmt.Printf("\033[1;31m║  %d SKIPS - LOCKING FOR 2 MINUTES!                        ║\033[0m\n", e.lockLimit)
	fmt.Printf("\033[1;31m╚══════════════════════════════════════════════════════════╝\033[0m\n")
	fmt.Println()

	// Lock loop - re-lock every 5 seconds for 2 minutes
	// Even if you unlock, it locks again
	for i := 1; i <= 24; i++ {
		logWarn(fmt.Sprintf("Lock %d/24 - wait it out!", i))
		lockScreen()
		time.Sleep(5 * time.Second)
	}
	logMsg("Lock loop complete. Take your breaks next time!")
}

func (e *Enforcer) doShutdown() {
	fmt.Println()
	fmt.Printf("\033[1;31m╔══════════════════════════════════════════════════════════╗\033[0m\n")
	fmt.Printf("\033[1;31m║  %d SKIPS - SHUTTING DOWN. GO OUTSIDE.                    ║\033[0m\n", e.shutdownLimit)
	fmt.Printf("\033[1;31m╚══════════════════════════════════════════════════════════╝\033[0m\n")
	fmt.Println()

	time.Sleep(2 * time.Second)
	shutdownComputer()
}

func (e *Enforcer) GetSkipCount() int {
	return e.skipCount
}
