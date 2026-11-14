# Configuration

The server options add fields to `natsserver.Options`.

## Types

| Type      | Description                                 |
| --------- | ------------------------------------------- |
| `Server`  | The embedded NATS server, with slog logging |
| `Options` | Adds `ReadyTimeout` to `natsserver.Options` |

## Options

| Field          | Type            | Description                                                          |
| -------------- | --------------- | -------------------------------------------------------------------- |
| `Options`      | `*nats.Options` | The standard NATS server options: host, port, authentication         |
| `ReadyTimeout` | `time.Duration` | The maximum time that `Start()` waits for the server to become ready |

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
