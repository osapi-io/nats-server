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

// Package server provides an embedded NATS server implementation.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// DefaultReadyTimeout is the time that Start waits for the server to accept
// connections when Options.ReadyTimeout is zero.
const DefaultReadyTimeout = 10 * time.Second

// New makes a Server from its logger and options. Start runs it. A nil
// logger is slog.Default().
func New(
	logger *slog.Logger,

	opts *Options,
) *Server {
	if logger == nil {
		logger = slog.Default()
	}

	return &Server{
		logger: logger,
		Opts:   opts,
	}
}

// readyPoll is the time that Start waits for readiness at each try, between
// checks of a fatal startup error, ctx and Stop. Tests change it.
var readyPoll = 50 * time.Millisecond

// Start runs the embedded NATS server and waits until it accepts
// connections, until ReadyTimeout passes, or until ctx ends or Stop runs.
// It refuses nil options, and a second Start before Stop. It attaches the
// logger before the server runs, with the Debug and Trace of the options.
// When the server logs a fatal error while it starts, such as a port that
// is in use, Start returns that error at once. If the server does not
// become ready, Start shuts it down before it returns the error. Start does
// not hold the lock while it waits, so Stop can end a Start that waits.
// ctx bounds only the wait: after Start returns, the server runs until
// Stop, also when ctx ends.
func (s *Server) Start(
	ctx context.Context,
) error {
	natsServer, timeout, startup, waitCtx, err := s.begin(ctx)
	if err != nil {
		return err
	}

	go natsServer.Start()

	waitErr := waitReady(waitCtx, natsServer, timeout, startup)

	s.mu.Lock()
	defer s.mu.Unlock()

	// A Stop that ran while the server became ready found no server to shut
	// down, so Start shuts it down here.
	if waitErr == nil {
		if err := waitCtx.Err(); err != nil {
			waitErr = fmt.Errorf("server start stopped: %w", err)
		}
	}

	s.cancelStart()
	s.cancelStart = nil
	if waitErr != nil {
		natsServer.Shutdown()
		natsServer.WaitForShutdown()
		s.starting = false

		return waitErr
	}

	s.natsServer = natsServer
	s.starting = false
	s.log().Info("nats server started successfully")

	return nil
}

// begin checks the state and the options, makes the server with its
// logger, and marks the Server as starting, all under the lock.
func (s *Server) begin(
	ctx context.Context,
) (NATSServerInstance, time.Duration, *startupLogger, context.Context, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.natsServer != nil || s.starting {
		return nil, 0, nil, nil, fmt.Errorf("server already started")
	}

	if s.Opts == nil || s.Opts.Options == nil {
		return nil, 0, nil, nil, fmt.Errorf("server options are nil")
	}

	timeout := s.Opts.ReadyTimeout
	if timeout < 0 {
		return nil, 0, nil, nil, fmt.Errorf("ready timeout %s is negative", timeout)
	}
	if timeout == 0 {
		timeout = DefaultReadyTimeout
	}

	natsServer, err := NewNATSServer(s.Opts.Options)
	if err != nil {
		return nil, 0, nil, nil, fmt.Errorf("error starting server: %w", err)
	}

	startup := &startupLogger{SlogWrapper: NewSlogWrapper(s.log())}
	natsServer.SetLogger(startup, s.Opts.Debug, s.Opts.Trace)

	waitCtx, cancel := context.WithCancel(ctx)
	s.cancelStart = cancel
	s.starting = true

	return natsServer, timeout, startup, waitCtx, nil
}

// waitReady waits until natsServer accepts connections. It returns the
// fatal error that the server logged, ctx's error, or a timeout error.
func waitReady(
	ctx context.Context,
	natsServer NATSServerInstance,
	timeout time.Duration,
	startup *startupLogger,
) error {
	deadline := time.Now().Add(timeout)

	for {
		if natsServer.ReadyForConnections(min(readyPoll, max(time.Until(deadline), 0))) {
			return nil
		}

		if cause := startup.fatal(); cause != "" {
			return fmt.Errorf("server did not start: %s", cause)
		}

		if err := ctx.Err(); err != nil {
			return fmt.Errorf("server start stopped: %w", err)
		}

		if !time.Now().Before(deadline) {
			if last := startup.lastError(); last != "" {
				return fmt.Errorf("server not ready for connections after %s: %s", timeout, last)
			}

			return fmt.Errorf("server not ready for connections after %s", timeout)
		}
	}
}

// Stop shuts down the embedded NATS server and waits until it has stopped.
// It also ends a Start that waits. A Stop with no running server does
// nothing. After Stop, Start can run the server again.
func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancelStart != nil {
		s.cancelStart()
	}

	if s.natsServer == nil {
		return
	}

	s.log().Info("shutting down nats server")
	s.natsServer.Shutdown()
	s.natsServer.WaitForShutdown()
	s.natsServer = nil
	s.log().Info("nats server shut down successfully")
}

// ClientURL is the URL that a client connects to, with the port that the
// server listens on, also when Options.Port is -1. It is empty when the
// server does not run.
func (s *Server) ClientURL() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.natsServer == nil {
		return ""
	}

	return s.natsServer.ClientURL()
}

// log is the Server's logger, or slog.Default() for a Server that New did
// not make.
func (s *Server) log() *slog.Logger {
	if s.logger == nil {
		return slog.Default()
	}

	return s.logger
}
