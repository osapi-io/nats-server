# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go package for running an embedded NATS server. Provides server lifecycle management, slog-based logging integration, and configurable options. Used by osapi-io projects (linked via `replace` in consuming project's `go.mod`).

## Development Commands

```bash
just fetch             # Fetch shared justfiles (run once or to update)
just deps              # Install all dependencies
just test              # Run all tests (lint + unit + coverage)
just go::unit          # Run unit tests only
just go::vet           # Run golangci-lint
just go::fmt           # Auto-format (gofumpt + golines)
just go::fmt-check     # Check formatting without modifying
just go::unit-cov      # Generate coverage report
go test -run TestName -v ./pkg/server/...  # Run a single test
```

## Package Structure

- **`pkg/server/`** - Core embedded NATS server library
  - `server.go` - Server struct, constructor, Start/Stop lifecycle
  - `logger.go` - SlogWrapper adapting slog.Logger to NATS Logger interface
  - `server_wrapper.go` - NewNATSServer factory (overridable for tests)
  - `types.go` - Shared types (Server, Options, NATSServerInstance)
  - `mocks/` - Generated mock implementations

## Code Standards (MANDATORY)

### Function Signatures

ALL function signatures MUST use multi-line format:
```go
func FunctionName(
    param1 type1,
    param2 type2,
) (returnType, error) {
}
```

### Testing

- Public tests: `*_public_test.go` in test package (`package server_test`) for exported functions
- Internal tests: `*_test.go` in same package (`package server`) for private functions
- Use `testify/suite` with table-driven patterns
- Use `golang/mock` for mocking interfaces

### Go Patterns

- Error wrapping: `fmt.Errorf("context: %w", err)`
- Early returns over nested if-else
- Unused parameters: rename to `_`
- Import order: stdlib, third-party, local (blank-line separated)

### Linting

golangci-lint with: errcheck, errname, goimports, govet, prealloc, predeclared, revive, staticcheck. Generated files (`*.gen.go`, `*.pb.go`) are excluded from formatting.

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/) with the
50/72 rule. Format: `type(scope): description`.

When committing via Claude Code, end with:
- `Co-Authored-By: Claude <noreply@anthropic.com>`
