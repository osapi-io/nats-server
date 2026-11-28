// Copyright (c) 2025 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package server

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// SlogWrapper sends the printf-style log calls of the NATS server to an
// slog.Logger. It formats a message only when the logger's handler takes its
// level. Make one with NewSlogWrapper.
type SlogWrapper struct {
	logger *slog.Logger
}

// NewSlogWrapper returns a SlogWrapper that logs to logger. A nil logger is
// slog.Default().
func NewSlogWrapper(
	logger *slog.Logger,
) *SlogWrapper {
	if logger == nil {
		logger = slog.Default()
	}

	return &SlogWrapper{logger: logger}
}

// log formats and writes one message at level, if the handler takes it.
func (l *SlogWrapper) log(
	level slog.Level,
	format string,
	v []interface{},
) {
	ctx := context.Background()
	if !l.logger.Enabled(ctx, level) {
		return
	}

	l.logger.Log(ctx, level, message(format, v))
}

// message is format with v. A format with no v is the message as it is, so
// that a "%" in it stays as it is.
func message(
	format string,
	v []interface{},
) string {
	if len(v) == 0 {
		return format
	}

	return fmt.Sprintf(format, v...)
}

// Noticef logs a formatted notice message at Info.
func (l *SlogWrapper) Noticef(
	format string,
	v ...interface{},
) {
	l.log(slog.LevelInfo, format, v)
}

// Warnf logs a formatted warning message at Warn.
func (l *SlogWrapper) Warnf(
	format string,
	v ...interface{},
) {
	l.log(slog.LevelWarn, format, v)
}

// Fatalf logs a formatted fatal message at Error. It does not stop the
// process.
func (l *SlogWrapper) Fatalf(
	format string,
	v ...interface{},
) {
	l.log(slog.LevelError, format, v)
}

// Errorf logs a formatted error message at Error.
func (l *SlogWrapper) Errorf(
	format string,
	v ...interface{},
) {
	l.log(slog.LevelError, format, v)
}

// Debugf logs a formatted debug message at Debug.
func (l *SlogWrapper) Debugf(
	format string,
	v ...interface{},
) {
	l.log(slog.LevelDebug, format, v)
}

// Tracef logs a formatted trace message at Debug, because slog has no trace
// level.
func (l *SlogWrapper) Tracef(
	format string,
	v ...interface{},
) {
	l.log(slog.LevelDebug, format, v)
}

// startupLogger is the SlogWrapper that Start gives the NATS server. It also
// keeps the first fatal message and the last error message, so that Start
// can return why the server did not start.
type startupLogger struct {
	*SlogWrapper

	mu        sync.Mutex
	firstFail string
	lastErr   string
}

// Fatalf logs a fatal message at Error, and keeps the first one.
func (l *startupLogger) Fatalf(
	format string,
	v ...interface{},
) {
	l.mu.Lock()
	if l.firstFail == "" {
		l.firstFail = message(format, v)
	}
	l.mu.Unlock()

	l.SlogWrapper.Fatalf(format, v...)
}

// Errorf logs an error message at Error, and keeps the last one.
func (l *startupLogger) Errorf(
	format string,
	v ...interface{},
) {
	l.mu.Lock()
	l.lastErr = message(format, v)
	l.mu.Unlock()

	l.SlogWrapper.Errorf(format, v...)
}

// fatal is the first fatal message, or empty.
func (l *startupLogger) fatal() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.firstFail
}

// lastError is the last error message, or empty.
func (l *startupLogger) lastError() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.lastErr
}
