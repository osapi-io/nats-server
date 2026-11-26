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

// Start runs the embedded NATS server and waits until it accepts
// connections. It refuses nil options, and a second Start before Stop. It
// attaches the logger before the server runs, with the Debug and Trace of the
// options, so the logger also gets the startup messages. If the server does
// not become ready, Start shuts it down before it returns the error.
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.natsServer != nil {
		return fmt.Errorf("server already started")
	}

	if s.Opts == nil || s.Opts.Options == nil {
		return fmt.Errorf("server options are nil")
	}

	timeout := s.Opts.ReadyTimeout
	if timeout < 0 {
		return fmt.Errorf("ready timeout %s is negative", timeout)
	}
	if timeout == 0 {
		timeout = DefaultReadyTimeout
	}

	natsServer, err := NewNATSServer(s.Opts.Options)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	natsServer.SetLogger(
		NewSlogWrapper(s.logger),
		s.Opts.Debug,
		s.Opts.Trace,
	)

	go natsServer.Start()

	if !natsServer.ReadyForConnections(timeout) {
		natsServer.Shutdown()
		natsServer.WaitForShutdown()

		return fmt.Errorf("server not ready for connections")
	}

	s.logger.Info("nats server started successfully")

	s.natsServer = natsServer

	return nil
}

// Stop shuts down the embedded NATS server and waits until it has stopped.
// A Stop with no running server does nothing. After Stop, Start can run the
// server again.
func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.natsServer == nil {
		return
	}

	s.logger.Info("shutting down nats server")
	s.natsServer.Shutdown()
	s.natsServer.WaitForShutdown()
	s.natsServer = nil
	s.logger.Info("nats server shut down successfully")
}
