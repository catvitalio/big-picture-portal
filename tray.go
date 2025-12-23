package main

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed assets/icon.ico
var iconData []byte

var (
	bpDisplayItems   map[string]*systray.MenuItem
	mainDisplayItems map[string]*systray.MenuItem

	bpAudioMenu    *systray.MenuItem
	mainAudioMenu  *systray.MenuItem
	bpAudioItems   []*systray.MenuItem
	mainAudioItems []*systray.MenuItem
	audioDevices   []AudioDevice
)

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Big Picture Portal")
	systray.SetTooltip("Big Picture Portal")

	bpDisplayItems = make(map[string]*systray.MenuItem)
	mainDisplayItems = make(map[string]*systray.MenuItem)

	bigPictureMenu := systray.AddMenuItem("Big Picture Display", "Select display for Big Picture mode")
	bpDisplayItems["internal"] = bigPictureMenu.AddSubMenuItem("Internal", "PC screen only")
	bpDisplayItems["external"] = bigPictureMenu.AddSubMenuItem("External", "Second screen only")
	bpDisplayItems["duplicate"] = bigPictureMenu.AddSubMenuItem("Duplicate", "Duplicate screens")
	bpDisplayItems["extend"] = bigPictureMenu.AddSubMenuItem("Extend", "Extend screens")

	mainMenu := systray.AddMenuItem("Main Display", "Select display for normal mode")
	mainDisplayItems["internal"] = mainMenu.AddSubMenuItem("Internal", "PC screen only")
	mainDisplayItems["external"] = mainMenu.AddSubMenuItem("External", "Second screen only")
	mainDisplayItems["duplicate"] = mainMenu.AddSubMenuItem("Duplicate", "Duplicate screens")
	mainDisplayItems["extend"] = mainMenu.AddSubMenuItem("Extend", "Extend screens")

	systray.AddSeparator()

	bpAudioMenu = systray.AddMenuItem("Big Picture Audio", "Select audio output for Big Picture mode")
	mainAudioMenu = systray.AddMenuItem("Main Audio", "Select audio output for normal mode")

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	loadAudioDevices()
	updateMenuState()
	go monitorBigPicture()
	go handleAudioMenuClicks()

	for mode, item := range bpDisplayItems {
		m := mode
		go func(menuItem *systray.MenuItem) {
			for range menuItem.ClickedCh {
				config.BigPictureDisplay = m
				saveConfig()
				updateMenuState()
			}
		}(item)
	}

	for mode, item := range mainDisplayItems {
		m := mode
		go func(menuItem *systray.MenuItem) {
			for range menuItem.ClickedCh {
				config.MainDisplay = m
				saveConfig()
				updateMenuState()
			}
		}(item)
	}

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func updateMenuState() {
	updateDisplayMenuState()
	updateAudioMenuState()
}

func updateDisplayMenuState() {
	for mode, item := range bpDisplayItems {
		if mode == config.BigPictureDisplay {
			item.Check()
		} else {
			item.Uncheck()
		}
	}

	for mode, item := range mainDisplayItems {
		if mode == config.MainDisplay {
			item.Check()
		} else {
			item.Uncheck()
		}
	}
}

func loadAudioDevices() {
	devices, err := getAudioDevices()
	if err != nil {
		return
	}
	audioDevices = devices

	bpAudioItems = nil
	mainAudioItems = nil

	for i := range audioDevices {
		device := &audioDevices[i]

		displayName := device.Name
		if device.DeviceState != 1 {
			displayName += " (Disabled)"
		}

		bpItem := bpAudioMenu.AddSubMenuItem(displayName, device.ID)
		mainItem := mainAudioMenu.AddSubMenuItem(displayName, device.ID)

		bpAudioItems = append(bpAudioItems, bpItem)
		mainAudioItems = append(mainAudioItems, mainItem)
	}

	updateAudioMenuState()
}

func updateAudioMenuState() {
	for i, item := range bpAudioItems {
		item.Uncheck()
		if audioDevices[i].ID == config.BigPictureAudio {
			item.Check()
		}
	}

	for i, item := range mainAudioItems {
		item.Uncheck()
		if audioDevices[i].ID == config.MainAudio {
			item.Check()
		}
	}
}

func handleAudioMenuClicks() {
	for i := range audioDevices {
		idx := i

		go func() {
			for {
				<-bpAudioItems[idx].ClickedCh
				config.BigPictureAudio = audioDevices[idx].ID
				saveConfig()
				updateAudioMenuState()
			}
		}()

		go func() {
			for {
				<-mainAudioItems[idx].ClickedCh
				config.MainAudio = audioDevices[idx].ID
				saveConfig()
				updateAudioMenuState()
			}
		}()
	}
}
