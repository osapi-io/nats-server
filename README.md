[![go report card](https://goreportcard.com/badge/github.com/osapi-io/nats-server?style=for-the-badge)](https://goreportcard.com/report/github.com/osapi-io/nats-server)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
![gitHub commit activity](https://img.shields.io/github/commit-activity/m/osapi-io/nats-server?style=for-the-badge)

# NATS Server

A Go package for running an embedded NATS server.

## Usage

See the [examples][] for more uses.

```golang
func main() {
	debug := true
	trace := debug
	logger := getLogger(debug)

	opts := server.NewOptions(
		server.WithDebug(debug),
		server.WithTrace(trace),
		server.WithStoreDir(".nats/jetstream/"),
		server.WithReadyTimeout(5*time.Second),
	)

	streamOpts := server.NewStreamOptions(
		server.WithStreamName("TASKS"),
		server.WithSubjects("tasks.*"),

		server.WithConsumer(server.NewConsumerOptions(
			server.WithDurable("osapi"),
			server.WithAckPolicy(nats.AckExplicitPolicy),
			server.WithMaxDeliver(5),
			server.WithAckWait(30*time.Second),
		)),
	)

	var sm server.Manager = server.New(logger, opts, streamOpts)
	err := sm.Start()
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	sm.Stop()
}
```

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
