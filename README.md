<p align="center">
  <picture>
    <source srcset="asset/logo-dark.svg" media="(prefers-color-scheme: dark)">
    <source srcset="asset/logo-light.svg" media="(prefers-color-scheme: light)">
    <img src="asset/logo-dark.svg" alt="nats-server" width="582">
  </picture>
</p>

<p align="center">A Go package for running an embedded NATS server.</p>

<p align="center">
  <a href="https://github.com/osapi-io/nats-server/releases/latest"><img alt="release" src="https://img.shields.io/github/release/osapi-io/nats-server.svg?style=for-the-badge"></a>
  <a href="https://codecov.io/gh/osapi-io/nats-server"><img alt="codecov" src="https://img.shields.io/codecov/c/github/osapi-io/nats-server?style=for-the-badge"></a>
  <a href="LICENSE"><img alt="license" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
  <a href="https://github.com/osapi-io/nats-server/actions/workflows/go.yml"><img alt="build" src="https://img.shields.io/github/actions/workflow/status/osapi-io/nats-server/go.yml?style=for-the-badge"></a>
  <a href="https://github.com/goreleaser"><img alt="powered by" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge"></a>
  <a href="https://conventionalcommits.org"><img alt="conventional commits" src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge"></a>
  <a href="https://nats.io"><img alt="nats" src="https://img.shields.io/badge/NATS-27AAE1?style=for-the-badge&logo=natsdotio&logoColor=white"></a>
  <a href="https://just.systems"><img alt="built with just" src="https://img.shields.io/badge/Built_with-Just-black?style=for-the-badge&logo=just&logoColor=white"></a>
  <img alt="gitHub commit activity" src="https://img.shields.io/github/commit-activity/m/osapi-io/nats-server?style=for-the-badge">
  <a href="https://pkg.go.dev/github.com/osapi-io/nats-server/pkg/server"><img alt="go reference" src="https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go&logoColor=white"></a>
</p>

<p align="center">
<b>A NATS server as a goroutine in your own process.</b>
</p>

<p align="center">
Run a NATS server inside your own process with JetStream support and slog
logging. No daemon to install and no broker to operate: the server is a
goroutine in the program that starts it.
</p>

## Install

```bash
go get github.com/osapi-io/nats-server
```

## Features

See the [server docs](docs/server/README.md) for quick start, authentication,
and per-feature reference.

| Feature              | Description                                               | Docs                                 | Source                              |
| -------------------- | --------------------------------------------------------- | ------------------------------------ | ----------------------------------- |
| Lifecycle management | Non-blocking `Start()` / graceful `Stop()` with readiness | [docs](docs/server/lifecycle.md)     | [`server.go`](pkg/server/server.go) |
| slog integration     | Adapts `slog.Logger` to the NATS server logging interface | [docs](docs/server/logging.md)       | [`logger.go`](pkg/server/logger.go) |
| Configuration        | Options for host, port, store dir, auth, and timeouts     | [docs](docs/server/configuration.md) | [`types.go`](pkg/server/types.go)   |

## Examples

Each example is a standalone Go program you can read and run.

| Example                                           | What it shows                         |
| ------------------------------------------------- | ------------------------------------- |
| [auth-none](examples/auth-none/main.go)           | Start a server without authentication |
| [auth-user-pass](examples/auth-user-pass/main.go) | Server with username/password auth    |
| [auth-nkeys](examples/auth-nkeys/main.go)         | Server with NKey authentication       |
| [simple-server](examples/simple-server/main.go)   | Minimal server startup and shutdown   |

## Documentation

- [Package documentation] on pkg.go.dev. API reference.

## Contributing

See the [Contributing](CONTRIBUTING.md) guide for prerequisites, setup,
conventions, and the PR workflow.

## License

The [MIT] License.

[mit]: LICENSE
[package documentation]: https://pkg.go.dev/github.com/osapi-io/nats-server/pkg/server
