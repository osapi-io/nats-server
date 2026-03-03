# Configuration

Server options extend `natsserver.Options` with additional fields.

## Types

| Type      | Description                                                     |
| --------- | --------------------------------------------------------------- |
| `Server`  | Embedded NATS server with slog logging                          |
| `Options` | Extends `natsserver.Options` with `ReadyTimeout`                |

## Options

| Field          | Type            | Description                                     |
| -------------- | --------------- | ----------------------------------------------- |
| `Options`      | `*nats.Options` | Standard NATS server options (host, port, auth)  |
| `ReadyTimeout` | `time.Duration` | Max wait time for server readiness after start   |

## Usage

```go
opts := &server.Options{
    Options: &natsserver.Options{
        Host:      "localhost",
        Port:      4222,
        JetStream: true,
        StoreDir:  ".nats/jetstream/",
    },
    ReadyTimeout: 5 * time.Second,
}

s := server.New(logger, opts)
```
