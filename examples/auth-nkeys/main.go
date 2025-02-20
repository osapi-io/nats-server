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

	nkeyUsers := []*natsserver.NkeyUser{
		{
			Nkey: "UAZMBGU3ASBL22E5WW6F3EFAW3CNUGGFRRYI6VRPPOHLNCNJZDTXFOPG", // Service 1
		},
		{
			Nkey: "UCL5D5YPOGDRFAS354LW7D3E4S7FOJXZGR3ULHGQFWTFQ3PCRTUAT3ZD", // Service 2 disabled
			Permissions: &natsserver.Permissions{
				Publish: &natsserver.SubjectPermission{
					Allow: []string{},
				},
				Subscribe: &natsserver.SubjectPermission{
					Allow: []string{},
				},
			},
		},
	}

	opts := &server.Options{
		Options: &natsserver.Options{
			Nkeys:     nkeyUsers,
			JetStream: true,
			Debug:     debug,
			Trace:     trace,
			StoreDir:  ".nats/jetstream/",
			NoSigs:    true,
			NoLog:     false,
		},
		ReadyTimeout: 5 * time.Second,
	}

	s := server.New(logger, opts)
	err := s.Start()
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	jsOpts := &client.ClientOptions{
		Host: s.Opts.Host,
		Port: s.Opts.Port,
		Auth: client.AuthOptions{
			AuthType: client.NKeyAuth,
			NKeyFile: ".nkeys/service1.seed",
		},
	}

	js, err := client.NewJetStreamContext(jsOpts)
	if err != nil {
		logger.Error("failed to create jetstream context", "error", err)
		os.Exit(1)
	}

	streamOpts := &client.StreamConfig{
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

	c := client.New(logger)
	if err := c.SetupJetStream(js, streamOpts); err != nil {
		logger.Error("failed setting up jetstream", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.Stop()
}
