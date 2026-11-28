// Copyright (c) 2026 John Dewey

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

package server_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/nats-server/pkg/server"
)

type LoggerPublicTestSuite struct {
	suite.Suite

	buf *bytes.Buffer
}

func (s *LoggerPublicTestSuite) SetupTest() {
	s.buf = &bytes.Buffer{}
}

func (s *LoggerPublicTestSuite) SetupSubTest() {
	s.SetupTest()
}

func (s *LoggerPublicTestSuite) wrapper(
	level slog.Level,
) *server.SlogWrapper {
	return server.NewSlogWrapper(slog.New(slog.NewTextHandler(
		s.buf,
		&slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}

				return a
			},
		},
	)))
}

func (s *LoggerPublicTestSuite) TestLogMethods() {
	tests := []struct {
		name    string
		level   slog.Level
		call    func(*server.SlogWrapper)
		wantLog string
	}{
		{
			name:    "Noticef logs at Info level",
			level:   slog.LevelDebug,
			call:    func(w *server.SlogWrapper) { w.Noticef("hello %s", "world") },
			wantLog: "level=INFO msg=\"hello world\"\n",
		},
		{
			name:    "a message with no args keeps a literal percent sign",
			level:   slog.LevelDebug,
			call:    func(w *server.SlogWrapper) { w.Noticef("100% ready") },
			wantLog: "level=INFO msg=\"100% ready\"\n",
		},
		{
			name:    "Warnf logs at Warn level",
			level:   slog.LevelDebug,
			call:    func(w *server.SlogWrapper) { w.Warnf("hello %s", "world") },
			wantLog: "level=WARN msg=\"hello world\"\n",
		},
		{
			name:    "Fatalf logs at Error level",
			level:   slog.LevelDebug,
			call:    func(w *server.SlogWrapper) { w.Fatalf("hello %s", "world") },
			wantLog: "level=ERROR msg=\"hello world\"\n",
		},
		{
			name:    "Errorf logs at Error level",
			level:   slog.LevelDebug,
			call:    func(w *server.SlogWrapper) { w.Errorf("hello %s", "world") },
			wantLog: "level=ERROR msg=\"hello world\"\n",
		},
		{
			name:    "Debugf logs at Debug level",
			level:   slog.LevelDebug,
			call:    func(w *server.SlogWrapper) { w.Debugf("hello %s", "world") },
			wantLog: "level=DEBUG msg=\"hello world\"\n",
		},
		{
			name:    "Tracef logs at Debug level",
			level:   slog.LevelDebug,
			call:    func(w *server.SlogWrapper) { w.Tracef("hello %s", "world") },
			wantLog: "level=DEBUG msg=\"hello world\"\n",
		},
		{
			name:    "a level that the handler does not take writes nothing",
			level:   slog.LevelInfo,
			call:    func(w *server.SlogWrapper) { w.Tracef("hello %s", "world") },
			wantLog: "",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.call(s.wrapper(tc.level))

			s.Equal(tc.wantLog, s.buf.String())
		})
	}
}

func (s *LoggerPublicTestSuite) TestNewSlogWrapper() {
	tests := []struct {
		name string
	}{
		{name: "a nil logger is slog.Default"},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			prev := slog.Default()
			defer slog.SetDefault(prev)
			slog.SetDefault(slog.New(slog.NewTextHandler(s.buf, nil)))

			server.NewSlogWrapper(nil).Noticef("hello")

			s.Contains(s.buf.String(), "msg=hello")
		})
	}
}

func TestLoggerPublicTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerPublicTestSuite))
}
