package main

import (
	"os/exec"
)

func shutdownComputer() {
	exec.Command("osascript", "-e", `tell app "System Events" to shut down`).Run()
}
