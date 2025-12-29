<div align="center">

# 🛡️ Enodia

**Per-App Network Access Controller for Windows**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails](https://img.shields.io/badge/Wails-2.11-blue?style=flat)](https://wails.io)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://react.dev)
[![Windows](https://img.shields.io/badge/Windows-10%2F11-0078D6?style=flat&logo=windows)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Early%20Development-orange.svg)]()

*Take control of which apps can access the internet on your Windows PC*

> ⚠️ **Early Development** — This app is under active development. Contributions and bug reports are welcome!

</div>

---

## ✨ Features

- **🔍 Auto-Discovery** — Automatically detects all installed Win32 and Microsoft Store (UWP) apps
- **🚫 One-Click Blocking** — Block any app's internet access with a single click
- **🔄 Persistent Rules** — Firewall rules survive reboots
- **⚡ Lightweight** — Native Windows app with minimal resource usage

## 📸 Screenshots

*Coming soon*

## 🚀 Quick Start

### Prerequisites

- Windows 10/11
- [Go 1.21+](https://golang.org/dl/)
- [Node.js 24+](https://nodejs.org/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Installation

```bash
# Clone the repository
git clone https://github.com/aman7ph/Enodia.git
cd Enodia

# Install Wails CLI (if not installed)
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Run in development mode
wails dev

# Build for production
wails build
```

### Running

> ⚠️ **Administrator privileges required** — Enodia needs admin rights to create Windows Firewall rules.

Right-click the executable and select "Run as administrator", or run from an elevated terminal.

## 🏗️ Project Structure

```
Enodia/
├── app.go                 # App lifecycle
├── wails_api.go           # Wails-exposed API methods
├── main.go                # Entry point
├── internal/
│   ├── apps/              # App discovery
│   │   ├── discovery.go   # Main entry
│   │   ├── win32.go       # Registry-based discovery
│   │   ├── store.go       # UWP/Store app discovery
│   │   ├── types.go       # InstalledApp struct
│   │   └── utils.go       # Helper functions
│   └── firewall/          # Windows Firewall management
│       ├── manager.go     # COM worker thread
│       ├── block.go       # Block/Unblock methods
│       ├── rules.go       # Rule creation
│       ├── state.go       # Get blocked apps
│       └── types.go       # Constants & types
└── frontend/              # React + Vite + shadcn/ui
    └── src/
        ├── App.tsx        # Main component
        └── components/    # UI components
```

## 🔧 How It Works

1. **Discovery** — Scans Windows Registry and queries `Get-AppxPackage` for installed apps
2. **Firewall Rules** — Creates Windows Firewall rules using COM API (`HNetCfg.FwPolicy2`)
3. **UWP Support** — Uses Package SID (App Container SID) for blocking Store apps
4. **Persistence** — Rules are stored by Windows Firewall and persist across reboots

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.21+ |
| Frontend | React 18 + TypeScript |
| UI Framework | [shadcn/ui](https://ui.shadcn.com/) + Tailwind CSS |
| Desktop | [Wails 2](https://wails.io/) |
| Firewall | Windows COM API (go-ole) |

## 📝 Roadmap

- [ ] System tray support
- [ ] App icons for Win32 apps
- [ ] Network traffic monitoring
- [ ] Scheduled blocking profiles
- [ ] Android support (future)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Wails](https://wails.io/) — Go + Web framework for desktop apps
- [shadcn/ui](https://ui.shadcn.com/) — Beautiful UI components
- [go-ole](https://github.com/go-ole/go-ole) — Go bindings for Windows COM

---

<div align="center">

**Made with ❤️ for Windows power users**

</div>
