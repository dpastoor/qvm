# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

QVM (Quarto Version Manager) is a command-line tool written in Go that manages multiple versions of Quarto, similar to how nvm manages Node.js versions or pyenv manages Python versions. It supports both global version management and per-shell session version switching.

## Development Commands

The project uses Just (justfile) as the build tool. Common commands:

- `just build` - Build the binary (outputs to ./qvm)
- `just test` - Run tests with coverage
- `just lint` - Run golangci-lint
- `just fmt` - Format code with gofumpt
- `just ci` - Run full CI pipeline (setup, build, test)
- `just run <args>` - Run the binary with arguments (e.g., `just run ls`)
- `just setup` - Install dependencies (go mod tidy)
- `just cover` - Open test coverage report

For testing specific patterns:
- `just test "" "" TestSpecificFunction`
- `just test "" ./internal/config/...`

## Architecture

### Core Structure
- **main.go** - Entry point, handles macOS XDG path normalization
- **cmd/** - CLI commands using Cobra framework
  - **root.go** - Main command setup and configuration
  - Individual command files (ls.go, install.go, use.go, etc.)
- **internal/** - Internal packages
  - **config/** - Configuration and filesystem operations
  - **gh/** - GitHub API client for fetching Quarto releases
  - **pipeline/** - Download and installation pipeline
  - **unarchive/** - Archive extraction utilities

### Key Components
- Uses Cobra for CLI and Viper for configuration
- GitHub API integration for fetching Quarto releases
- XDG Base Directory specification for config storage (normalized on macOS)
- Interactive prompts using AlecAivazis/survey
- Parallel installation support
- Archive handling for multiple formats

### Configuration
- Follows XDG specification but normalizes macOS paths to use `.config` and `.local/share`
- Global configuration stored in XDG config directories
- Version management through filesystem links and PATH manipulation

## Testing and Quality

- Use `just test` for running tests with race detection and coverage
- Use `just lint` to run golangci-lint
- Use `just fmt` to format code with gofumpt
- All CI steps can be run locally with `just ci`

## Release Process

- `just release` - Creates and pushes new git tag using svu
- Uses GoReleaser for building multi-platform binaries
- Supports packaging for rpm, deb, and apk formats