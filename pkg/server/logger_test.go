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

package server

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testHandler struct {
	records []slog.Record
}

func (h *testHandler) Enabled(
	_ context.Context,
	_ slog.Level,
) bool {
	return true
}

func (h *testHandler) Handle(
	_ context.Context,
	r slog.Record,
) error {
	h.records = append(h.records, r)
	return nil
}

func (h *testHandler) WithAttrs(
	_ []slog.Attr,
) slog.Handler {
	return h
}

func (h *testHandler) WithGroup(
	_ string,
) slog.Handler {
	return h
}

type SlogWrapperTestSuite struct {
	suite.Suite

	handler *testHandler
	wrapper *SlogWrapper
}

func (s *SlogWrapperTestSuite) SetupTest() {
	s.handler = &testHandler{}
	s.wrapper = &SlogWrapper{
		logger: slog.New(s.handler),
	}
}

func (s *SlogWrapperTestSuite) SetupSubTest() {
	s.SetupTest()
}

func (s *SlogWrapperTestSuite) TestLogMethods() {
	tests := []struct {
		name          string
		call          func()
		expectedLevel slog.Level
	}{
		{
			name: "Noticef logs at Info level",
			call: func() {
				s.wrapper.Noticef("hello %s", "world")
			},
			expectedLevel: slog.LevelInfo,
		},
		{
			name: "Warnf logs at Warn level",
			call: func() {
				s.wrapper.Warnf("hello %s", "world")
			},
			expectedLevel: slog.LevelWarn,
		},
		{
			name: "Fatalf logs at Error level",
			call: func() {
				s.wrapper.Fatalf("hello %s", "world")
			},
			expectedLevel: slog.LevelError,
		},
		{
			name: "Errorf logs at Error level",
			call: func() {
				s.wrapper.Errorf("hello %s", "world")
			},
			expectedLevel: slog.LevelError,
		},
		{
			name: "Debugf logs at Debug level",
			call: func() {
				s.wrapper.Debugf("hello %s", "world")
			},
			expectedLevel: slog.LevelDebug,
		},
		{
			name: "Tracef logs at Debug level",
			call: func() {
				s.wrapper.Tracef("hello %s", "world")
			},
			expectedLevel: slog.LevelDebug,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.call()

			s.Require().Len(s.handler.records, 1)
			s.Equal(tc.expectedLevel, s.handler.records[0].Level)
			s.Equal("hello world", s.handler.records[0].Message)
		})
	}
}

func TestSlogWrapperTestSuite(t *testing.T) {
	suite.Run(t, new(SlogWrapperTestSuite))
}
