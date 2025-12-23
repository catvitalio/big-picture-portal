# Big Picture Portal

A Windows system tray application that automatically switches display and audio output configurations when entering or exiting Steam Big Picture mode.

## Features

- 🖥️ **Automatic Display Switching** - Seamlessly switch between display configurations (Internal, External, Duplicate, Extend) when launching Steam Big Picture
- 🔊 **Automatic Audio Switching** - Automatically switch audio output devices (e.g., from PC speakers to TV via HDMI)
- 🎮 **Steam Big Picture Detection** - Monitors Steam Big Picture mode and applies your preferred settings
- 🔄 **Reversible** - Returns to your main configuration when exiting Big Picture mode
- ⚙️ **Configurable** - Easy-to-use system tray interface for setting preferences
- 💾 **Persistent Settings** - Remembers your configuration across restarts

## Use Case

Perfect for users who:
- Play Steam games on a TV via HDMI while using their PC monitor for regular work
- Want different audio outputs for gaming vs regular use
- Need automatic switching without manual intervention

## Requirements

- Windows 10/11
- Go 1.21+ (for building from source)
- Steam (for Big Picture mode detection)

## Installation

### Download Pre-built Binary

Download the latest release from the [Releases](../../releases) page and run `BigPicturePortal.exe`.

### Build from Source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/big-picture-portal.git
cd big-picture-portal
```

2. Install dependencies:
```bash
go get
```

3. Build using the provided batch file:
```bash
build.bat
```

This will create `big-picture-portal.exe` in the project directory.

Alternatively, build manually:
```bash
go build -ldflags="-H windowsgui" -o big-picture-portal.exe .
```

The `-H windowsgui` flag prevents a console window from appearing.

## Usage

1. Run `big-picture-portal.exe`
2. The application will appear in your system tray
3. Right-click the tray icon to configure settings:
   - **Big Picture Display** - Choose display mode for Big Picture (Internal/External/Duplicate/Extend)
   - **Main Display** - Choose display mode for normal use
   - **Big Picture Audio** - Select audio output for Big Picture (optional)
   - **Main Audio** - Select audio output for normal use (optional)
4. Launch Steam Big Picture mode - settings will automatically apply
5. Exit Big Picture mode - settings will automatically revert

### Configuration Options

#### Display Modes
- **Internal** - PC screen only
- **External** - Second screen only (e.g., TV)
- **Duplicate** - Mirror displays
- **Extend** - Extended desktop across both screens

#### Audio Devices
- Lists all available audio output devices
- Devices marked as "(Disabled)" are currently inactive (e.g., HDMI audio when TV is off)
- You can pre-select disabled devices - they'll be activated when available

## Configuration File

Settings are stored in:
```
%APPDATA%\BigPicturePortal\config.json
```

Example configuration:
```json
{
  "bigPictureDisplay": "external",
  "mainDisplay": "internal",
  "bigPictureAudio": "{device-id-here}",
  "mainAudio": "{device-id-here}",
  "checkInterval": 2000
}
```

- `checkInterval` - Milliseconds between Big Picture mode checks (default: 2000)

## How It Works

1. **Monitoring** - Polls every 2 seconds for Steam Big Picture mode window
2. **Detection** - Detects both English and localized window titles
3. **Display Switching** - Uses Windows `DisplaySwitch.exe` API
4. **Audio Switching** - Uses Windows Core Audio API via COM interfaces

## Technical Details

### Dependencies

- [github.com/getlantern/systray](https://github.com/getlantern/systray) - System tray functionality
- [github.com/go-ole/go-ole](https://github.com/go-ole/go-ole) - COM interface support
- [github.com/moutend/go-wca](https://github.com/moutend/go-wca) - Windows Core Audio API
- [golang.org/x/sys/windows](https://pkg.go.dev/golang.org/x/sys/windows) - Windows system calls

### Architecture

```
├── main.go       - Entry point
├── tray.go       - System tray UI and menu handling
├── steam.go      - Steam Big Picture detection and monitoring
├── display.go    - Display switching logic
├── audio.go      - Audio device enumeration and switching
├── config.go     - Configuration management
└── assets/
    └── icon.ico  - System tray icon
```

## Building

### Prerequisites

```bash
go version  # Ensure Go 1.21+
```

### Build Commands

**Standard build:**
```bash
go build -o big-picture-portal.exe .
```

**Release build (no console window):**
```bash
go build -ldflags="-H windowsgui" -o big-picture-portal.exe .
```

**Using build.bat:**
```bash
build.bat
```

The batch file builds with the windowsgui flag automatically.


## License

MIT License - see LICENSE file for details

## Credits

Built with Go and Windows APIs for seamless gaming experience.

Icon provided by [Glaze Icon Pack](https://grafikstash.com/glaze/) from GrafikStash.
