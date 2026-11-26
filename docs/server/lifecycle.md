# Lifecycle management

`Start()` returns when the server is ready for connections. `Stop()` closes the
connections of the server and stops it.

## Methods

| Method    | Description                                                 |
| --------- | ----------------------------------------------------------- |
| `Start()` | Starts the embedded NATS server and waits until it is ready |
| `Stop()`  | Closes the connections and stops the NATS server            |

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
is ready for connections. It waits no longer than `ReadyTimeout`. When
`ReadyTimeout` is zero, it waits `DefaultReadyTimeout`, which is 10 seconds. A
negative `ReadyTimeout` is an error. If the server is not ready in that time,
`Start()` shuts the server down, waits until it stops, and returns an error.
Thus a failed `Start()` leaves no listener and no goroutine. When the server is
ready, `Start()` sets the slog logger on it.

`Start()` returns an error if the server already runs. Call `Stop()` first.

`Stop()` calls `Shutdown()` on the NATS server and waits until the server stops.
When `Stop()` returns, the port and the JetStream store are free. A second
`Stop()` does nothing, and a `Start()` after `Stop()` runs the server again.
`Start()` and `Stop()` are safe to call from different goroutines.
