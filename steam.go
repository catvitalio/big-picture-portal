package main

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func isSteamBigPictureRunning() bool {
	titles := []string{
		"Режим Big Picture",
		"Steam Big Picture Mode",
	}

	for _, title := range titles {
		hwnd := findWindowByTitle(title)
		if hwnd != 0 {
			return true
		}
	}

	return false
}

func findWindowByTitle(title string) uintptr {
	user32 := windows.NewLazySystemDLL("user32.dll")
	findWindowW := user32.NewProc("FindWindowW")

	titlePtr, _ := windows.UTF16PtrFromString(title)

	ret, _, _ := findWindowW.Call(
		0,
		uintptr(unsafe.Pointer(titlePtr)),
	)

	return ret
}

func monitorBigPicture() {
	ticker := time.NewTicker(time.Duration(config.CheckInterval) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		bigPictureRunning := isSteamBigPictureRunning()

		if bigPictureRunning {
			switchDisplay(config.BigPictureDisplay)
		} else if !bigPictureRunning {
			switchDisplay(config.MainDisplay)
		}
	}
}
