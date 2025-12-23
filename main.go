package main

import "github.com/getlantern/systray"

func main() {
	if err := initConfig(); err != nil {
		return
	}

	systray.Run(onReady, func() {})
}
