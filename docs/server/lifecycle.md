# Lifecycle management

`Start()` returns when the server is ready for connections. `Stop()` shuts down
the server gracefully.

## Methods

| Method    | Description                                                 |
| --------- | ----------------------------------------------------------- |
| `Start()` | Starts the embedded NATS server and waits until it is ready |
| `Stop()`  | Shuts down the NATS server gracefully                       |

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

`Start()` starts the NATS server in a goroutine. Then it waits until the server
is ready for connections. It waits no longer than `ReadyTimeout`. If the server
is not ready in that time, `Start()` returns an error. When the server is ready,
`Start()` sets the slog logger on it.

`Stop()` calls `Shutdown()` on the NATS server. The NATS server then closes its
connections and stops.
