[![release](https://img.shields.io/github/release/osapi-io/nats-server.svg?style=for-the-badge)](https://github.com/osapi-io/nats-server/releases/latest)
[![go report card](https://goreportcard.com/badge/github.com/osapi-io/nats-server?style=for-the-badge)](https://goreportcard.com/report/github.com/osapi-io/nats-server)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/osapi-io/nats-server/go.yml?style=for-the-badge)](https://github.com/osapi-io/nats-server/actions/workflows/go.yml)
[![powered by](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
![gitHub commit activity](https://img.shields.io/github/commit-activity/m/osapi-io/nats-server?style=for-the-badge)

# NATS Server

A Go package for running an embedded NATS server.

## Usage

https://github.com/osapi-io/nats-server/blob/b542d33de18737037c3e4c6ed8160b27440690fb/examples/auth-none/main.go#L21-L63

See the [examples][] section for additional use cases.

## Documentation

See the [generated documentation][] for details on available packages and functions.

## Development

Fetch shared justfiles:

```bash
just fetch
```

Install dependencies:

```bash
just deps
```

Run all tests:

```bash
just test
```

Auto format code:

```bash
just go::fmt
```

List available recipes:

```bash
just --list
```

## License

The [MIT][] License.

[examples]: examples/
[generated documentation]: docs/gen/
[MIT]: LICENSE
