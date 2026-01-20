# 🦙 ALPAKA NOTES v2.0 - Ultimate TUI Edition

**The most beautiful** note-taking application for your terminal! Built with ❤️ using Charm Bracelet Bubble Tea.

```
   ▄▄▄       ██▓     ██▓███   ▄▄▄       ██ ▄█▀▄▄▄      
  ▒████▄    ▓██▒    ▓██░  ██▒▒████▄     ██▄█▒▒████▄    
  ▒██  ▀█▄  ▒██░    ▓██░ ██▓▒▒██  ▀█▄  ▓███▄░▒██  ▀█▄  
  ░██▄▄▄▄██ ▒██░    ▒██▄█▓▒ ▒░██▄▄▄▄██ ▓██ █▄░██▄▄▄▄██ 
   ▓█   ▓██▒░██████▒▒██▒ ░  ░ ▓█   ▓██▒▒██▒ █▄▓█   ▓██▒
```

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](https://github.com)

## ✨ Premium Features

### 🎨 **Beautiful Interface**
- Gradient colors and animations
- ASCII art logo
- Animated splash screen
- Responsive layout
- Colorful tags and icons
- Blinking cursor animations

### 🔐 **Security**
- Password-protected encryption
- SHA-256 password hashing
- Custom `.alpaka` file format
- No password storage
- Show/hide password toggle (Ctrl+H)

### 📝 **Note Features**
- Unlimited notes
- Multi-line content support
- Tag system
- Character counters
- Automatic timestamps
- Real-time search

### 📊 **Statistics & Analytics**
- Note, word, and tag counts
- Average words per note
- Tag cloud visualization
- Recent activity tracking
- Usage charts

### 🎯 **3 View Modes**
1. **List** - Detailed preview of all notes
2. **Grid** - Compact 2-column view
3. **Preview** - Full view of single note

### 🔄 **3 Sorting Modes**
- By date (newest first)
- By title (alphabetically)
- By tags

## 🚀 Quick Start

### Installation

#### Option 1: Download Pre-built Binary (Easiest)

Download the latest release for your platform:
- **Linux**: `alpaka-linux-amd64`
- **Windows**: `alpaka-windows-amd64.exe`
- **macOS Intel**: `alpaka-macos-intel`
- **macOS Apple Silicon**: `alpaka-macos-m1`

```bash
# Linux/macOS
chmod +x alpaka-linux-amd64
./alpaka-linux-amd64

# Windows
# Double-click alpaka-windows-amd64.exe
```

#### Option 2: Build from Source

**Requirements:**
- Go 1.18+ (recommended 1.21+)

```bash
# Clone or create directory
mkdir alpaka-notes && cd alpaka-notes

# Copy files: main.go, screens.go, notebook.go, go.mod

# Install dependencies
go mod tidy

# Run!
go run .
```

## 🎮 Controls

### Global
- **Ctrl+C** - Exit application
- **Esc** - Return to main menu
- **↑/↓** or **j/k** - Navigate (Vim keys!)

### Login Screen
- Type password
- **Ctrl+H** - Show/hide password
- **Enter** - Login

### Main Menu
- **↑/↓** or **j/k** - Select option
- **Enter** - Confirm
- **q** - Quit

### Add Note
- **Tab** - Next field
- **Shift+Tab** - Previous field
- **Enter** - New line (in content)
- **Ctrl+S** - Save note
- **Esc** - Cancel

### Browse Notes
- **↑/↓** or **j/k** - Scroll
- **d** - Delete note
- **v** - Change view (List/Grid/Preview)
- **s** - Change sorting
- **Esc** - Return

### Search
- Type query
- Real-time results
- **Esc** - Return

### Statistics
- Browse data
- **Esc** - Return

### Settings
- **↑/↓** - Select option
- **Enter/Space** - Change setting
- **Esc** - Return

## 📁 Project Structure

```
alpaka-notes/
├── main.go          # Main application + styles
├── screens.go       # All screens (Login, Menu, etc.)
├── notebook.go      # Data model + encryption
├── go.mod           # Dependencies
├── go.sum           # Checksums
├── README.md        # This documentation
└── notatki.alpaka   # Your encrypted notes
```

## 🎨 Color Palette

```
Pink:    #FF6B9D  - Primary
Purple:  #C792EA  - Secondary
Blue:    #82AAFF  - Accent
Cyan:    #89DDFF  - Info
Green:   #C3E88D  - Success
Yellow:  #FFCB6B  - Warning
Orange:  #F78C6C  - Highlight
Red:     #FF5370  - Danger
```

## 📦 .alpaka File Format

```
ALPAKA
VERSION:1.0
HASH:<sha256_password_hash>
---ENCRYPTED---
<encrypted_json_data>
```

## 🔨 Building

### Cross-platform compilation

```bash
# All platforms
chmod +x build.sh
./build.sh

# Or use Makefile
make all-platforms

# Individual platforms
make linux
make windows
make macos
```

### Manual build

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o alpaka-linux-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o alpaka-windows-amd64.exe

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o alpaka-macos-intel

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o alpaka-macos-m1
```

## 🔒 Security

**Current implementation:**
- XOR cipher (demonstration)
- SHA-256 password hash
- JSON serialization

**Planned improvements:**
- [ ] AES-256-GCM encryption
- [ ] PBKDF2/Argon2 key derivation
- [ ] Salt generation
- [ ] Encrypted backups

## 🎯 Roadmap

### v2.1
- [ ] Export to Markdown/PDF
- [ ] Import from other formats
- [ ] Categories/folders
- [ ] Pinned notes
- [ ] Archive

### v2.2
- [ ] Note attachments
- [ ] Inline images
- [ ] Markdown rendering
- [ ] Syntax highlighting
- [ ] Color themes

### v3.0
- [ ] Cloud sync (optional)
- [ ] Note sharing
- [ ] Web interface
- [ ] Mobile app
- [ ] Plugins

## 🏆 Premium Features Included

✅ **Animated splash screen** with gradients  
✅ **7 screens** (Splash, Login, Menu, Add, Browse, Search, Stats, Settings)  
✅ **ASCII art logo** with gradient  
✅ **Blinking cursor** in all input fields  
✅ **Colorful tags** (5 colors rotating)  
✅ **3 view modes** for notes  
✅ **3 sorting modes**  
✅ **Statistics** with counters and charts  
✅ **Real-time search**  
✅ **Character counters** in forms  
✅ **Progress bar** on loading  
✅ **Status messages** (Success/Error/Warning/Info)  
✅ **Vim keybindings** (j/k)  
✅ **Responsive** layout  

## 📚 Libraries

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI Framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style & Layout

## 🌟 Screenshots

```
╔═══════════════════════════════════════════════╗
║          🦙 ALPAKA NOTES v2.0 🦙             ║
║       Ultimate TUI Experience Edition        ║
╚═══════════════════════════════════════════════╝

📊 Statistics: 42 notes | File: notatki.alpaka
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ▶ 📝  New Note          - Create new entry
    📖  Browse             - View all notes
    🔍  Search             - Find notes
    📊  Statistics         - Analytics and charts
    ⚙️   Settings          - Sorting and views
    💾  Save               - Save changes to disk
    🚪  Exit               - Close program
```

## 🤝 Contributing

Want to help? Great!

1. Fork the project
2. Create a branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open a Pull Request

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

- [Charm Bracelet](https://charm.sh/) - For amazing TUI tools
- All contributors and users
- The Go community

## 📞 Support

- 🐛 **Bug reports**: [Open an issue](https://github.com/yourusername/alpaka-notes/issues)
- 💡 **Feature requests**: [Open an issue](https://github.com/yourusername/alpaka-notes/issues)
- ❓ **Questions**: [Discussions](https://github.com/yourusername/alpaka-notes/discussions)

## ⭐ Star History

If you like this project, please give it a star! ⭐

---

**Made with ❤️ and 🦙 by Alpaka Inc.**

*Secure • Beautiful • Fast*
