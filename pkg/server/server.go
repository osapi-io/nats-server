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
	"github.com/nats-io/nats.go"
)

// New initialize and configure a new Server instance.
func New(
	logger *slog.Logger,
	options *Options,
	streamOptions ...*StreamOptions,
) *Server {
	server := &Server{
		logger:  logger,
		options: options,
	}

	if options.JetStream && len(streamOptions) > 0 {
		server.streamOptions = streamOptions
	}

	return server
}

// Start start the embedded NATS server.
func (s *Server) Start() error {
	natsOpts := s.options.ToNATSOptions()

	natsServer, err := natsserver.NewServer(natsOpts)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	go natsServer.Start()

	// Wait for server readiness
	if !natsServer.ReadyForConnections(s.options.ReadyTimeout) {
		return fmt.Errorf("server not ready for connections")
	}

	slogWrapper := &SlogWrapper{
		logger: s.logger,
	}

	natsServer.SetLogger(slogWrapper, true, true)

	s.logger.Info("nats server started successfully")

	if s.options.JetStream && len(s.streamOptions) > 0 {
		if err := s.setupJetStream(); err != nil {
			return fmt.Errorf("error setting up jetstream: %w", err)
		}

		s.logger.Info("jet stream setup completed successfully")
	}

	s.natsServer = natsServer

	return nil
}

// setupJetStream creates the JetStream connection and stream configuration.
func (s *Server) setupJetStream() error {
	if len(s.streamOptions) == 0 {
		return fmt.Errorf("jetstream is enabled but no stream configuration was provided")
	}

	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", s.options.Host, s.options.Port))
	if err != nil {
		return fmt.Errorf("error connecting to server: %w", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("error enabling jetstream: %w", err)
	}

	for _, stream := range s.streamOptions {
		natsStreamConfig := stream.ToNATS()

		_, err := js.AddStream(natsStreamConfig)
		if err != nil {
			return fmt.Errorf("error creating stream %s: %w", stream.Name, err)
		}

		// Iterate over each consumer tied to the stream
		for _, consumer := range stream.Consumers {
			natsConsumerConfig := consumer.ToNATS()
			_, err := js.AddConsumer(stream.Name, natsConsumerConfig)
			if err != nil {
				return fmt.Errorf("error creating consumer for stream %s: %w", stream.Name, err)
			}
		}
	}

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
