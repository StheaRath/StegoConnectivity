# StegoConnectivity

Desktop steganography tool built with Go + Wails + React. Backend logic lives in Go (`internal/stego`, `internal/crypto`); frontend is Vite/React and embedded into the binary at build time.

## Prerequisites
- Go 1.24+
- Node.js 16+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- macOS builds: Xcode Command Line Tools

## Install (once)
```bash
cd frontend
npm install
```

## Development (hot reload)
```bash
wails dev
```

## Build
- mac universal: `wails build -clean -platform darwin/universal`
- mac per-arch: `wails build -clean -platform darwin/arm64` or `darwin/amd64`
- Windows x86_64: `wails build -clean -platform windows/amd64` (add `-nsis` for installer if `makensis` is on PATH)

Outputs land in `build/bin/` (mac .app / Windows .exe).

## Running binaries
- mac: run the .app in `build/bin/` (or move to /Applications)
- Windows: run the .exe in `build/bin/` (WebView2 runtime required; usually present on Win11). If using an installer, build with `-nsis`.

## Key files
- App entry/options: `main.go`, `app.go`, `wails.json`
- Frontend: `frontend/src`, config `frontend/vite.config.js`
- Backend logic: `internal/stego`, `internal/crypto`
