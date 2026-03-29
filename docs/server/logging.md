# slog Integration

Adapts Go's `slog.Logger` to the NATS server `Logger` interface, bridging
structured logging with NATS's printf-style log calls.

## Type

| Type          | Description                                     |
| ------------- | ----------------------------------------------- |
| `SlogWrapper` | Implements NATS `Logger` interface using `slog` |

## Methods

| Method      | Maps to        |
| ----------- | -------------- |
| `Noticef()` | `slog.Info()`  |
| `Warnf()`   | `slog.Warn()`  |
| `Fatalf()`  | `slog.Error()` |
| `Errorf()`  | `slog.Error()` |
| `Debugf()`  | `slog.Debug()` |
| `Tracef()`  | `slog.Debug()` |

The wrapper uses `fmt.Sprintf` to format NATS's printf-style messages before
passing them to slog's structured logging methods.
