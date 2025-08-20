# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Velo is a cross-platform framework for creating mobile applications using Go with UI from web technologies (React, Vue, Svelte, Angular). It uses a WebView approach to wrap web applications into native Android, iOS, desktop, and web applications.

## Development Commands

### Core CLI Commands
- `go run cmd/velo/main.go <command>` - Run the Velo CLI
- `make init` - Initialize project (shortcut for `go run cmd/velo/main.go init`)

### Available CLI Commands
- `velo init` - Initialize a new Velo project with framework selection
- `velo build` - Build the application for production
- `velo dev` - Start development server with hot reload
- `velo doctor` - Check system requirements and dependencies
- `velo version` - Show version information
- `velo help` - Display help information

### Frontend Development
- Frontend projects use Vite as the build tool
- Main frontend package.json scripts:
  - `npm run dev` - Start Vite dev server
  - `npm run build` - Build frontend for production
  - `npm run preview` - Preview production build

### React App (cmd/velo/sss/)
- Uses Create React App
- Scripts: `start`, `build`, `test`, `eject`

## Architecture

### Project Structure
```
velo/
├── cmd/velo/           # CLI entry point
├── pkg/                # Core packages
│   ├── cli/           # CLI implementation
│   ├── server/        # Preview server
│   ├── builder/       # Platform builders
│   └── utils/         # Utilities
├── internal/          # Internal packages
│   ├── desktop/       # Desktop platform code
│   ├── android/       # Android platform code
│   ├── ios/          # iOS platform code
│   └── web/          # Web platform code
├── constants/         # Framework and command definitions
├── platform/         # Platform-specific assets and builds
├── frontend/         # Main web application
└── examples/         # Example applications
```

### Key Components

#### CLI System (`pkg/cli/`)
- Uses a command registry pattern with constants for command definitions
- Commands are implemented in `pkg/cli/commands/`
- Main CLI entry point in `cmd/velo/main.go`

#### Framework Support (`internal/framework.go`, `constants/framework.go`)
Supports multiple frontend frameworks:
- **React**: Create React App, Next.js, Remix, Vite
- **Vue**: Nuxt, Quasar, Vite
- **Svelte**: SvelteKit, Vite
- **Angular**: Universal, Nest.js
- **Solid**: SolidStart, Vite
- **Astro**: Vite

#### Platform Builders (`pkg/builder/`)
- `android.go` - Android app building
- `ios.go` - iOS app building  
- `frontend.go` - Frontend bundling

#### Preview Server (`pkg/server/`)
- HTTP file server for previewing built assets
- Serves from `mobile-shell/assets/` directory

### Platform Architecture
- **Desktop**: Uses WebView with GTK4/WebKit (Linux), native webviews (Windows/macOS)
- **Mobile**: WebView wrappers for Android and iOS
- **Web**: Standard web deployment
- **Bridge Communication**: Two-way communication between Go backend and JavaScript frontend

## Dependencies

### Go Dependencies
- `github.com/charmbracelet/huh` - Terminal UI forms
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/webview/webview_go` - Cross-platform webview

### Frontend Dependencies
- Vite for build tooling
- React ecosystem support
- Expo for React Native capabilities

## Development Notes

### Working with Platforms
- Desktop Linux implementation uses CGO with WebKit
- Platform-specific code is organized under `internal/` by platform
- Assets are bundled to `platform/assets/` for mobile deployment

### Command Implementation
- New commands should be added to `constants/command.go`
- Command implementations go in `pkg/cli/commands/`
- CLI routing is handled in `pkg/cli/cli.go`

### Framework Installation
- Framework installers are in `internal/framework.go`
- Uses npm/npx for web framework scaffolding
- Interactive prompts supported for framework-specific configuration

### Testing
- No existing test files found - tests should be added as development progresses
- Frontend projects include testing frameworks (Jest for React apps)