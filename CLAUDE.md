# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`projj-go` is a Go rewrite of [projj](https://github.com/popomore/projj), a tool for managing code repositories organized by host/owner/repo structure (e.g., `~/Code/github.com/user/repo`).

## Commands

Once `go.mod` is initialized, standard Go commands apply:

```bash
# Build
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./path/to/package -run TestName

# Lint (requires golangci-lint)
golangci-lint run

# Format
gofmt -w .
# or
goimports -w .
```

## Notes

- This project is in early development — no source code exists yet beyond the `.gitignore`.
- The `.gitignore` is configured for Go (excludes binaries, `*.test`, `coverage.*`, `go.work`).
