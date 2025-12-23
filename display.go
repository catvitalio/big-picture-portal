package main

import (
	"os/exec"
)

func switchDisplay(displayType string) {
	var arg string
	switch displayType {
	case "internal":
		arg = "/internal"
	case "external":
		arg = "/external"
	case "extend":
		arg = "/extend"
	case "duplicate":
		arg = "/clone"
	default:
		return
	}

	cmd := exec.Command("DisplaySwitch.exe", arg)
	cmd.Run()
}
