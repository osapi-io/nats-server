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
	"fmt"
	"log/slog"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

// New initialize and configure a new Server instance.
func New(
	logger *slog.Logger,

	opts *Options,
) *Server {
	return &Server{
		logger: logger,
		Opts:   opts,
	}
}

// Start start the embedded NATS server.
func (s *Server) Start() error {
	natsServer, err := natsserver.NewServer(s.Opts.Options)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	go natsServer.Start()

	// Wait for server readiness
	if !natsServer.ReadyForConnections(s.Opts.ReadyTimeout) {
		return fmt.Errorf("server not ready for connections")
	}

	slogWrapper := &SlogWrapper{
		logger: s.logger,
	}

	natsServer.SetLogger(slogWrapper, true, true)

	s.logger.Info("nats server started successfully")

	s.natsServer = natsServer

	return nil
}

// Stop gracefully stops the embedded NATS server.
func (s *Server) Stop() {
	if s.natsServer != nil {
		s.logger.Info("shutting down nats server")
		s.natsServer.Shutdown()
		s.logger.Info("nats server shut down successfully")
	}
}
