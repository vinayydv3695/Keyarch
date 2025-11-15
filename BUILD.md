# Keyarch Build and Installation Guide

## Prerequisites

- Go 1.22 or higher
- Git
- A terminal with Unicode and color support

## Building from Source

### 1. Clone the Repository

```bash
git clone https://github.com/vinayydv3695/keyarch.git
cd keyarch
```

### 2. Install Dependencies

```bash
go mod tidy
```

This will download all required dependencies:
- github.com/charmbracelet/bubbletea
- github.com/charmbracelet/lipgloss
- github.com/charmbracelet/bubbles
- github.com/spf13/cobra
- modernc.org/sqlite

### 3. Build the Application

#### For your current platform:

```bash
go build -o keyarch ./cmd/keyarch
```

#### With optimizations (smaller binary):

```bash
go build -ldflags="-s -w" -o keyarch ./cmd/keyarch
```

### 4. Install (Optional)

#### Linux/macOS:

```bash
sudo mv keyarch /usr/local/bin/
```

Or to your local bin:

```bash
mkdir -p ~/.local/bin
mv keyarch ~/.local/bin/
# Add ~/.local/bin to your PATH if not already
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc  # or ~/.zshrc
```

#### Windows:

Move `keyarch.exe` to a directory in your PATH or add its directory to PATH.

## Cross-Compilation

### For Linux (amd64):

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o keyarch-linux-amd64 ./cmd/keyarch
```

### For Linux (arm64):

```bash
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o keyarch-linux-arm64 ./cmd/keyarch
```

### For macOS (Intel):

```bash
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o keyarch-darwin-amd64 ./cmd/keyarch
```

### For macOS (Apple Silicon):

```bash
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o keyarch-darwin-arm64 ./cmd/keyarch
```

### For Windows:

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o keyarch-windows-amd64.exe ./cmd/keyarch
```

## Quick Install via Go

If you have Go installed, you can install directly:

```bash
go install github.com/vinayydv3695/keyarch/cmd/keyarch@latest
```

This will install the binary to `$GOPATH/bin` (usually `~/go/bin`).

## Running the Application

### TUI Mode (Interactive):

```bash
keyarch
```

### Direct Mode Examples:

```bash
# Quick 15-second test
keyarch --mode quick

# 30-second test
keyarch --mode quick --duration 30

# 50-word test
keyarch --mode words --words 50

# Quote mode
keyarch --mode quote

# Code practice mode
keyarch --mode code

# With specific theme
keyarch --theme dracula
```

## Troubleshooting

### Issue: "command not found: keyarch"

**Solution**: The binary is not in your PATH. Either:
- Add the directory containing `keyarch` to your PATH
- Move `keyarch` to a directory already in PATH (like `/usr/local/bin`)
- Run it with full path: `./keyarch`

### Issue: Terminal doesn't display colors properly

**Solution**: 
- Make sure your terminal supports 256 colors or TrueColor
- Try a modern terminal: Alacritty, Kitty, iTerm2, Windows Terminal
- Set your TERM variable: `export TERM=xterm-256color`

### Issue: Unicode characters not displaying

**Solution**:
- Ensure your terminal uses UTF-8 encoding
- Install a font with good Unicode support (Nerd Fonts, JetBrains Mono, etc.)

### Issue: Database errors

**Solution**:
- Remove the database: `rm ~/.keyarch/keyarch.db`
- The app will create a fresh database on next run

### Issue: Build fails with dependency errors

**Solution**:
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod tidy

# Try building again
go build ./cmd/keyarch
```

## Development Setup

### Running in Development:

```bash
go run ./cmd/keyarch
```

### Running with hot reload (using air):

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Create .air.toml config
air init

# Run
air
```

### Running Tests:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/engine/
```

## File Locations

After first run, Keyarch creates:

- **Config**: `~/.keyarch/config.json`
- **Database**: `~/.keyarch/keyarch.db`

To reset everything:

```bash
rm -rf ~/.keyarch
```

## Performance Tips

### Reduce Binary Size:

```bash
# Use UPX compression (install upx first)
go build -ldflags="-s -w" -o keyarch ./cmd/keyarch
upx --best --lzma keyarch
```

### Build with Race Detector (for testing):

```bash
go build -race -o keyarch ./cmd/keyarch
```

## Distribution

### Creating a Release:

```bash
# Create release directory
mkdir -p release

# Build for all platforms
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o release/keyarch-linux-amd64 ./cmd/keyarch
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o release/keyarch-darwin-amd64 ./cmd/keyarch
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o release/keyarch-darwin-arm64 ./cmd/keyarch
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o release/keyarch-windows-amd64.exe ./cmd/keyarch

# Create archives
cd release
tar -czf keyarch-linux-amd64.tar.gz keyarch-linux-amd64
tar -czf keyarch-darwin-amd64.tar.gz keyarch-darwin-amd64
tar -czf keyarch-darwin-arm64.tar.gz keyarch-darwin-arm64
zip keyarch-windows-amd64.zip keyarch-windows-amd64.exe
```

## Docker (Optional)

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o keyarch ./cmd/keyarch

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/keyarch /usr/local/bin/
ENTRYPOINT ["keyarch"]
```

Build and run:

```bash
docker build -t keyarch .
docker run -it --rm keyarch
```

---

For more information, visit: https://github.com/vinayydv3695/keyarch
