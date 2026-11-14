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
it gives the message to the slog method in the table.
