# NATS server

The `server` package gives an embedded NATS server. The server has JetStream and
options that you can configure. It sends its logs to slog. Make a server with
`New()`. Call `Start()` to run it.

## Quick start

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

if err := s.Start(context.Background()); err != nil {
    log.Fatal(err)
}
defer s.Stop()
```

## Features

| Feature                             | Description                                            | Source      |
| ----------------------------------- | ------------------------------------------------------ | ----------- |
| [`Lifecycle`](lifecycle.md)         | Start and Stop, with a readiness check                 | `server.go` |
| [`Logs`](logging.md)                | Connects slog to the NATS `Logger` interface           | `logger.go` |
| [`Configuration`](configuration.md) | An `Options` struct that adds to nats-server's options | `types.go`  |

## Authentication

Give the authentication options in `natsserver.Options`:

| Mode                  | Description              | Options field |
| --------------------- | ------------------------ | ------------- |
| No authentication     | Accept all connections   | (default)     |
| Username and password | List of user credentials | `Users`       |
| NKey                  | List of public NKeys     | `Nkeys`       |
