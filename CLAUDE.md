# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go package for running an embedded NATS server. Provides server lifecycle management, slog-based logging integration, and configurable options. Used by osapi-io projects (linked via `replace` in consuming project's `go.mod`).

## Development Reference

For setup, prerequisites, and contributing guidelines:

- @docs/development.md - Prerequisites, setup, code style, testing, commits
- @docs/contributing.md - PR workflow and contribution guidelines

## Quick Reference

```bash
just fetch / just deps / just test / just go::unit / just go::vet / just go::fmt
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

### Branching

See @docs/development.md#branching for full conventions.

When committing changes via `/commit`, create a feature branch first if
currently on `main`. Branch names use the pattern `type/short-description`
(e.g., `feat/add-dns-retry`, `fix/memory-leak`, `docs/update-readme`).

### Commit Messages

See @docs/development.md#commit-messages for full conventions.

Follow [Conventional Commits](https://www.conventionalcommits.org/) with the
50/72 rule. Format: `type(scope): description`.

When committing via Claude Code, end with:
- `🤖 Generated with [Claude Code](https://claude.ai/code)`
- `Co-Authored-By: Claude <noreply@anthropic.com>`
