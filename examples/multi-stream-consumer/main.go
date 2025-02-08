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
	"github.com/nats-io/nats.go"
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

	opts := server.NewOptions(
		server.WithDebug(debug),
		server.WithTrace(trace),
		server.WithStoreDir(".nats/jetstream/"),
		server.WithReadyTimeout(5*time.Second),
	)

	streamOpts1 := server.NewStreamOptions(
		server.WithStreamName("TASK_QUEUE"),
		server.WithSubjects("tasks.*"),
		server.WithConsumer(server.NewConsumerOptions(
			server.WithDurable("worker1"),
			server.WithAckPolicy(nats.AckExplicitPolicy),
			server.WithDeliverPolicy(nats.DeliverNewPolicy),
			server.WithMaxAckPending(10),
			server.WithAckWait(30*time.Second),
		)),
		server.WithConsumer(server.NewConsumerOptions(
			server.WithDurable("worker2"),
			server.WithAckPolicy(nats.AckExplicitPolicy),
			server.WithDeliverPolicy(nats.DeliverNewPolicy),
			server.WithMaxAckPending(10),
			server.WithAckWait(30*time.Second),
		)),
	)

	streamOpts2 := server.NewStreamOptions(
		server.WithStreamName("STREAM2"),
		server.WithSubjects("stream2.*"),
		server.WithConsumer(server.NewConsumerOptions(
			server.WithDurable("consumer3"),
			server.WithAckPolicy(nats.AckExplicitPolicy),
			server.WithMaxDeliver(5),
			server.WithAckWait(30*time.Second),
		)),
		server.WithConsumer(server.NewConsumerOptions(
			server.WithDurable("consumer4"),
			server.WithAckPolicy(nats.AckExplicitPolicy),
			server.WithMaxDeliver(5),
			server.WithAckWait(30*time.Second),
		)),
	)

	var sm server.Manager = server.New(logger, opts, streamOpts1, streamOpts2)
	err := sm.Start()
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	sm.Stop()
}
