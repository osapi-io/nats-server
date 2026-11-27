# slog integration

`SlogWrapper` connects Go's `slog.Logger` to the `Logger` interface of the NATS
server. The NATS server makes printf-style log calls. The wrapper sends them to
slog, which makes structured logs.

## Type

| Type          | Description                                        |
| ------------- | -------------------------------------------------- |
| `SlogWrapper` | Implements the NATS `Logger` interface with `slog` |

## Methods

| Method      | Maps to        |
| ----------- | -------------- |
| `Noticef()` | `slog.Info()`  |
| `Warnf()`   | `slog.Warn()`  |
| `Fatalf()`  | `slog.Error()` |
| `Errorf()`  | `slog.Error()` |
| `Debugf()`  | `slog.Debug()` |
| `Tracef()`  | `slog.Debug()` |

The wrapper formats each printf-style message from NATS with `fmt.Sprintf`. Then
it gives the message to the slog method in the table. It formats a message only
when the handler takes its level. `Fatalf()` logs. It does not stop the process.

`NewSlogWrapper(logger)` makes a wrapper. A nil logger is `slog.Default()`.
`Start()` makes one from the logger that `New` got.

## Debug and trace

NATS writes debug and trace messages only when `Options.Debug` and
`Options.Trace` are true. Both are false by default. Trace writes each protocol
line, so turn it on only to find a fault.

`slog` has no trace level. Thus a trace message goes to `slog.Debug`, and the
handler can't tell it from a debug message. Trace shows only when the handler
takes Debug.

```go
s := server.New(logger, &server.Options{
    Options: &natsserver.Options{Debug: true, Trace: true},
})
```
