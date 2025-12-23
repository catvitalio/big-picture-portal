package main

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed assets/icon.ico
var iconData []byte

var (
	// Big Picture display items
	bpInternalItem  *systray.MenuItem
	bpExternalItem  *systray.MenuItem
	bpExtendItem    *systray.MenuItem
	bpDuplicateItem *systray.MenuItem

	// Main display items
	mainInternalItem  *systray.MenuItem
	mainExternalItem  *systray.MenuItem
	mainExtendItem    *systray.MenuItem
	mainDuplicateItem *systray.MenuItem
)

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Big Picture Portal")
	systray.SetTooltip("Big Picture Portal")

	// Big Picture display menu
	bigPictureMenu := systray.AddMenuItem("Big Picture Display", "Select display for Big Picture mode")
	bpInternalItem = bigPictureMenu.AddSubMenuItem("Internal", "PC screen only")
	bpExternalItem = bigPictureMenu.AddSubMenuItem("External", "Second screen only")
	bpDuplicateItem = bigPictureMenu.AddSubMenuItem("Duplicate", "Duplicate screens")
	bpExtendItem = bigPictureMenu.AddSubMenuItem("Extend", "Extend screens")

	systray.AddSeparator()

	// Main display menu
	mainMenu := systray.AddMenuItem("Main Display", "Select display for normal mode")
	mainInternalItem = mainMenu.AddSubMenuItem("Internal", "PC screen only")
	mainExternalItem = mainMenu.AddSubMenuItem("External", "Second screen only")
	mainDuplicateItem = mainMenu.AddSubMenuItem("Duplicate", "Duplicate screens")
	mainExtendItem = mainMenu.AddSubMenuItem("Extend", "Extend screens")

	systray.AddSeparator()

	// Quit menu
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Update menu checkmarks
	updateMenuState()

	// Start monitoring
	go monitorBigPicture()

	// Handle menu clicks
	go func() {
		for {
			select {
			// Big Picture display options
			case <-bpInternalItem.ClickedCh:
				config.BigPictureDisplay = "internal"
				saveConfig()
				updateMenuState()

			case <-bpExternalItem.ClickedCh:
				config.BigPictureDisplay = "external"
				saveConfig()
				updateMenuState()

			case <-bpDuplicateItem.ClickedCh:
				config.BigPictureDisplay = "duplicate"
				saveConfig()
				updateMenuState()

			case <-bpExtendItem.ClickedCh:
				config.BigPictureDisplay = "extend"
				saveConfig()
				updateMenuState()

			// Main display options
			case <-mainInternalItem.ClickedCh:
				config.MainDisplay = "internal"
				saveConfig()
				updateMenuState()

			case <-mainExternalItem.ClickedCh:
				config.MainDisplay = "external"
				saveConfig()
				updateMenuState()

			case <-mainDuplicateItem.ClickedCh:
				config.MainDisplay = "duplicate"
				saveConfig()
				updateMenuState()

			case <-mainExtendItem.ClickedCh:
				config.MainDisplay = "extend"
				saveConfig()
				updateMenuState()

			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func updateMenuState() {
	// Update Big Picture display checkmarks
	bpInternalItem.Uncheck()
	bpExternalItem.Uncheck()
	bpDuplicateItem.Uncheck()
	bpExtendItem.Uncheck()

	switch config.BigPictureDisplay {
	case "internal":
		bpInternalItem.Check()
	case "external":
		bpExternalItem.Check()
	case "duplicate":
		bpDuplicateItem.Check()
	case "extend":
		bpExtendItem.Check()
	}

	// Update Main display checkmarks
	mainInternalItem.Uncheck()
	mainExternalItem.Uncheck()
	mainDuplicateItem.Uncheck()
	mainExtendItem.Uncheck()

	switch config.MainDisplay {
	case "internal":
		mainInternalItem.Check()
	case "external":
		mainExternalItem.Check()
	case "duplicate":
		mainDuplicateItem.Check()
	case "extend":
		mainExtendItem.Check()
	}
}
