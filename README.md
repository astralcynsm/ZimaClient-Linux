# ZimaClient

A Linux desktop client for managing ZimaOS devices, built with Wails (Go + Vue 3 + Naive UI).

## Features

- **Smart Multi-Path Connection**: Automatically switches between LAN and ZeroTier networks
- **Device Discovery**: Auto-scan for Zima devices on local network
- **Storage Management**: Mount and manage SMB shares with auto-mount wizard
- **ZeroTier Integration**: Built-in ZeroTier management and auto-installation
- **Modern UI**: Clean, responsive interface with light/dark mode support

## Dependencies

### Runtime Dependencies

- **WebKitGTK 4.0** (required)
  ```bash
  # Debian/Ubuntu
  sudo apt install libwebkit2gtk-4.0-37

  # Arch/EndeavourOS
  sudo pacman -S webkit2gtk

  # Fedora
  sudo dnf install webkit2gtk3
  ```

- **SMB Client** (for storage mounting)
  ```bash
  # Debian/Ubuntu
  sudo apt install cifs-utils smbclient

  # Arch/EndeavourOS
  sudo pacman -S cifs-utils smbclient

  # Fedora
  sudo dnf install cifs-utils samba-client
  ```

- **ZeroTier** (auto-installed by the app if not present)

- **Optional**: File dialog tools
  - `zenity` (GNOME) or `kdialog` (KDE) for folder selection dialogs

### Build Dependencies

- **Go** 1.21 or higher
- **Node.js** 18.x or higher
- **Wails CLI** v2
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

- **Development libraries**:
  ```bash
  # Debian/Ubuntu
  sudo apt install build-essential libgtk-3-dev libwebkit2gtk-4.0-dev

  # Arch/EndeavourOS
  sudo pacman -S base-devel gtk3 webkit2gtk

  # Fedora
  sudo dnf install @development-tools gtk3-devel webkit2gtk3-devel
  ```

## Installation

### Using AppImage (Recommended)

1. Download `ZimaClient-x86_64.AppImage` from releases
2. Make it executable:
   ```bash
   chmod +x ZimaClient-x86_64.AppImage
   ```
3. Run:
   ```bash
   ./ZimaClient-x86_64.AppImage
   ```

**Note**: WebKitGTK 4.0 must be installed on your system (not bundled to keep AppImage size small).

### Building from Source

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd rnctl
   ```

2. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   cd ..
   ```

3. Build:
   ```bash
   wails build
   ```

4. The binary will be in `build/bin/`

### Building AppImage

Use the provided Docker build script:
```bash
./docker-build.sh
```

This builds the AppImage in an Ubuntu 20.04 container for maximum compatibility.

## Important Notes

### Mount Points

- Default mount location: `/run/user/$(id -u)/zimaclient/<device-share>`
- Uses runtime directory to avoid file manager performance issues
- Mounts are automatically cleaned up on logout

### Sudo/Polkit

- The app requires elevated privileges for mounting operations
- On first mount, you'll be prompted for your sudo password
- Password is stored encrypted in memory for the session
- Systems with `pkexec` can use PolicyKit instead

### ZeroTier

- ZeroTier is auto-detected and installed if missing
- Embedded packages: `.deb` (3.5MB) and `.rpm` (12MB)
- Falls back to online installation if embedded packages fail
- Requires sudo privileges for installation

### Compatibility

- **Tested on**: EndeavourOS (Arch-based), Ubuntu 20.04+
- **Should work on**: Most modern Linux distributions with WebKitGTK 4.0
- **Desktop Environments**: GNOME, KDE, XFCE, and others
- **Architecture**: x86_64 only

### Known Limitations

- WebKitGTK 4.0 is not bundled in AppImage (system dependency)
- Requires sudo password or PolicyKit for mount operations
- ZeroTier installation requires internet connection if embedded packages fail

## Configuration

Settings are stored in:
- User preferences: `localStorage` (browser-like storage)
- Encrypted credentials: `~/.config/rnctl/vault/`
- Auto-mount config: Encrypted in vault

## Troubleshooting

### "Failed to load WebKit"
Install WebKitGTK 4.0 for your distribution (see Runtime Dependencies).

### "Mount failed"
- Ensure `cifs-utils` and `smbclient` are installed
- Check if you have sudo privileges
- Verify the device is accessible on the network

### "ZeroTier installation failed"
- Check internet connection
- Manually install ZeroTier: `curl -s https://install.zerotier.com | sudo bash`

### Device not found during scan
- Ensure device and computer are on the same network
- Check if device is powered on and network cable is connected
- Try using ZeroTier network ID for remote connection

## Development

### Project Structure

```
rnctl/
├── app.go              # Main Wails app backend
├── main.go             # Entry point
├── pkg/                # Go packages
│   ├── auth/           # Authentication & sudo management
│   ├── network/        # Smart connection manager
│   ├── storage/        # SMB mounting & management
│   ├── utils/          # System utilities
│   └── zerotier/       # ZeroTier integration & scanning
└── frontend/           # Vue 3 frontend
    ├── src/
    │   ├── views/      # Page components
    │   ├── components/ # Reusable components
    │   └── router/     # Vue Router config
    └── wailsjs/        # Auto-generated Wails bindings
```

### Running in Dev Mode

```bash
wails dev
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
