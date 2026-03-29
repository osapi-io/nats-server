# Lifecycle Management

Non-blocking server startup with graceful shutdown.

## Methods

| Method    | Description                                        |
| --------- | -------------------------------------------------- |
| `Start()` | Start the embedded NATS server, wait for readiness |
| `Stop()`  | Gracefully shut down the NATS server               |

## Usage

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

`Start()` launches the NATS server in a goroutine, waits for it to be ready for
connections (up to `ReadyTimeout`), then configures slog-based logging. `Stop()`
calls `Shutdown()` on the underlying NATS server for graceful cleanup.
