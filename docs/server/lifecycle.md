# Lifecycle management

`Start(ctx)` returns when the server is ready for connections. `Stop()` closes
the connections of the server and stops it. `ClientURL()` is the URL that a
client connects to.

## Methods

| Method        | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| `Start(ctx)`  | Starts the embedded NATS server and waits until it is ready  |
| `Stop()`      | Closes the connections and stops the NATS server             |
| `ClientURL()` | The URL to connect to, with the port it listens on, or empty |

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

if err := s.Start(context.Background()); err != nil {
    log.Fatal(err)
}
defer s.Stop()
```

`Start(ctx)` starts the NATS server in a goroutine. Then it waits until the
server is ready for connections. It waits no longer than `ReadyTimeout`. When
`ReadyTimeout` is zero, it waits `DefaultReadyTimeout`, which is 10 seconds. A
negative `ReadyTimeout` is an error.

The wait stops early in three cases:

- The server logs a fatal message while it starts, such as a port that is in
  use. `Start` returns `server did not start:` and that message at once.
- ctx ends. `Start` returns `server start stopped:` and ctx's error.
- `Stop()` runs. `Start` does not hold the lock while it waits, so `Stop()` does
  not block, and `Start` returns `server start stopped`.

When `ReadyTimeout` passes, the error names the timeout and the last error that
the server logged. In each case, `Start` shuts the server down and waits until
it stops. Thus a failed `Start` leaves no listener and no goroutine.

`Start` sets the slog logger on the server before the server runs. Thus the
startup messages go to the logger. `Start` returns an error if the options or
the NATS options are nil.

`Start` returns an error if the server already runs, or while a start waits.
Call `Stop()` first.

With `Port: -1`, the server picks a free port. `ClientURL()` returns the URL
with that port, such as `nats://127.0.0.1:52811`.

`Stop()` calls `Shutdown()` on the NATS server and waits until the server stops.
When `Stop()` returns, the port and the JetStream store are free. A second
`Stop()` does nothing, and a `Start()` after `Stop()` runs the server again.
`Start()` and `Stop()` are safe to call from different goroutines.
