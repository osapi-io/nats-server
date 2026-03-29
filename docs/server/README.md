# NATS Server

The `server` package provides an embedded NATS server with JetStream support,
slog-based logging, and configurable options. Create a server with `New()` and
call `Start()` to run it.

## Quick Start

```go
s := server.New(logger, &server.Options{
    Options: &natsserver.Options{
        Host:      "localhost",
        Port:      4222,
        JetStream: true,
        StoreDir:  ".nats/jetstream/",
    },
    ReadyTimeout: 5 * time.Second,
})

if err := s.Start(); err != nil {
    log.Fatal(err)
}
defer s.Stop()
```

## Features

| Feature                             | Description                                  | Source      |
| ----------------------------------- | -------------------------------------------- | ----------- |
| [`Lifecycle`](lifecycle.md)         | Non-blocking Start/Stop with readiness check | `server.go` |
| [`Logging`](logging.md)             | slog adapter for the NATS Logger interface   | `logger.go` |
| [`Configuration`](configuration.md) | Options struct extending nats-server options | `types.go`  |

## Authentication

Pass authentication options through `natsserver.Options`:

| Type          | Description              | Options field |
| ------------- | ------------------------ | ------------- |
| No auth       | Accept all connections   | (default)     |
| Username/pass | List of user credentials | `Users`       |
| NKey          | List of public NKeys     | `Nkeys`       |
