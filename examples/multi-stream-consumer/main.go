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

package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/osapi-io/nats-client/pkg/client"
	"github.com/osapi-io/nats-server/pkg/server"
)

func getLogger(debug bool) *slog.Logger {
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.Kitchen,
		}),
	)

	return logger
}

func main() {
	debug := true
	trace := debug
	logger := getLogger(debug)

	opts := &server.Options{
		Options: &natsserver.Options{
			JetStream: true,
			Debug:     debug,
			Trace:     trace,
			StoreDir:  ".nats/jetstream/",
			NoSigs:    true,
			NoLog:     false,
		},
		ReadyTimeout: 5 * time.Second,
	}

	streamOpts1 := &client.StreamConfig{
		StreamConfig: &nats.StreamConfig{
			Name:     "TASK_QUEUE",
			Subjects: []string{"tasks.*"},
			Storage:  nats.FileStorage,
			Replicas: 1,
		},
		Consumers: []*client.ConsumerConfig{
			{
				ConsumerConfig: &nats.ConsumerConfig{
					Durable:       "worker1",
					AckPolicy:     nats.AckExplicitPolicy,
					MaxAckPending: 10,
					AckWait:       30 * time.Second,
				},
			},
			{
				ConsumerConfig: &nats.ConsumerConfig{
					Durable:       "worker2",
					AckPolicy:     nats.AckExplicitPolicy,
					MaxAckPending: 10,
					AckWait:       30 * time.Second,
				},
			},
		},
	}

	streamOpts2 := &client.StreamConfig{
		StreamConfig: &nats.StreamConfig{
			Name:     "STREAM2",
			Subjects: []string{"stream2.*"},
			Storage:  nats.FileStorage,
			Replicas: 1,
		},
		Consumers: []*client.ConsumerConfig{
			{
				ConsumerConfig: &nats.ConsumerConfig{
					Durable:    "consumer3",
					AckPolicy:  nats.AckExplicitPolicy,
					MaxDeliver: 5,
					AckWait:    30 * time.Second,
				},
			},
			{
				ConsumerConfig: &nats.ConsumerConfig{
					Durable:    "consumer4",
					AckPolicy:  nats.AckExplicitPolicy,
					MaxDeliver: 5,
					AckWait:    30 * time.Second,
				},
			},
		},
	}

	s := server.New(logger, opts)
	err := s.Start()
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	js, err := client.NewJetStreamContext(s.Opts.Host, s.Opts.Port)
	if err != nil {
		logger.Error("failed to create jetstream context", "error", err)
		os.Exit(1)
	}

	c := client.New(logger, streamOpts1, streamOpts2)
	if err := c.SetupJetStream(js); err != nil {
		logger.Error("failed setting up jetstream", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.Stop()
}
