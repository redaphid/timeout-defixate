package main

import (
	"bufio"
	"os/exec"
	"strings"
)

type Event int

const (
	EventPostponed Event = iota
	EventSkipped
	EventCompleted
)

func (e Event) String() string {
	switch e {
	case EventPostponed:
		return "Postponed"
	case EventSkipped:
		return "Skipped"
	case EventCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

// StartMonitor spawns `log stream` and returns a channel of events
func StartMonitor() (<-chan Event, error) {
	events := make(chan Event)

	cmd := exec.Command("log", "stream",
		"--predicate", `process == "Time Out"`,
		"--info",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		defer close(events)
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			line := scanner.Text()

			// Skip lines starting with - (timestamp separators)
			if strings.HasPrefix(line, "-") {
				continue
			}

			if strings.Contains(line, "Postponed the break") {
				events <- EventPostponed
			} else if strings.Contains(line, "Skipped the break") {
				events <- EventSkipped
			} else if strings.Contains(line, "finished: yes") {
				events <- EventCompleted
			}
		}
	}()

	return events, nil
}
