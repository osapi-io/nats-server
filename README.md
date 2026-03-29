[![release](https://img.shields.io/github/release/osapi-io/nats-server.svg?style=for-the-badge)](https://github.com/osapi-io/nats-server/releases/latest)
[![codecov](https://img.shields.io/codecov/c/github/osapi-io/nats-server?style=for-the-badge)](https://codecov.io/gh/osapi-io/nats-server)
[![go report card](https://goreportcard.com/badge/github.com/osapi-io/nats-server?style=for-the-badge)](https://goreportcard.com/report/github.com/osapi-io/nats-server)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/osapi-io/nats-server/go.yml?style=for-the-badge)](https://github.com/osapi-io/nats-server/actions/workflows/go.yml)
[![powered by](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![nats](https://img.shields.io/badge/NATS-27AAE1?style=for-the-badge&logo=natsdotio&logoColor=white)](https://nats.io)
[![built with just](https://img.shields.io/badge/Built_with-Just-black?style=for-the-badge&logo=just&logoColor=white)](https://just.systems)
![gitHub commit activity](https://img.shields.io/github/commit-activity/m/osapi-io/nats-server?style=for-the-badge)
[![go reference](https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/osapi-io/nats-server/pkg/server)

# NATS Server

A Go package for running an embedded NATS server.

## 📦 Install

```bash
go get github.com/osapi-io/nats-server
```

## ✨ Features

See the [server docs](docs/server/README.md) for quick start, authentication,
and per-feature reference.

| Feature              | Description                                               | Docs                                       | Source                                    |
| -------------------- | --------------------------------------------------------- | ------------------------------------------ | ----------------------------------------- |
| Lifecycle management | Non-blocking `Start()` / graceful `Stop()` with readiness | [docs](docs/server/lifecycle.md)           | [`server.go`](pkg/server/server.go)       |
| slog integration     | Adapts `slog.Logger` to the NATS server logging interface | [docs](docs/server/logging.md)             | [`logger.go`](pkg/server/logger.go)       |
| Configuration        | Options for host, port, store dir, auth, and timeouts     | [docs](docs/server/configuration.md)       | [`types.go`](pkg/server/types.go)         |

## 📋 Examples

Each example is a standalone Go program you can read and run.

| Example                                          | What it shows                            |
| ------------------------------------------------ | ---------------------------------------- |
| [auth-none](examples/auth-none/main.go)         | Start a server without authentication    |
| [auth-user-pass](examples/auth-user-pass/main.go) | Server with username/password auth     |
| [auth-nkeys](examples/auth-nkeys/main.go)       | Server with NKey authentication          |
| [simple-server](examples/simple-server/main.go) | Minimal server startup and shutdown      |

## 📖 Documentation

See the [package documentation][] on pkg.go.dev for API details.

## 🤝 Contributing

See the [Development](docs/development.md) guide for prerequisites, setup,
and conventions. See the [Contributing](docs/contributing.md) guide before
submitting a PR.

## 📄 License

The [MIT][] License.

[package documentation]: https://pkg.go.dev/github.com/osapi-io/nats-server/pkg/server
[MIT]: LICENSE
