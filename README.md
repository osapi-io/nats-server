[![go report card](https://goreportcard.com/badge/github.com/osapi-io/nats-server?style=for-the-badge)](https://goreportcard.com/report/github.com/osapi-io/nats-server)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
![gitHub commit activity](https://img.shields.io/github/commit-activity/m/osapi-io/nats-server?style=for-the-badge)

# NATS Server

A Go package for running an embedded NATS server.

## Usage

https://github.com/osapi-io/nats-server/blob/b542d33de18737037c3e4c6ed8160b27440690fb/examples/auth-none/main.go#L21-L63

See the [examples][] section for additional use cases.

## Testing

Enable [Remote Taskfile][] feature.

```bash
export TASK_X_REMOTE_TASKFILES=1
```

Install dependencies:

```bash
$ task go:deps
```

To execute tests:

```bash
$ task go:test
```

Auto format code:

```bash
$ task go:fmt
```

List helpful targets:

```bash
$ task
```

## License

The [MIT][] License.

[examples]: examples/
[Remote Taskfile]: https://taskfile.dev/experiments/remote-taskfiles/
[MIT]: LICENSE
