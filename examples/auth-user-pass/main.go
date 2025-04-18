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

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/osapi-io/nats-server/pkg/server"
)

func main() {
	debug := true
	trace := debug
	logger := slog.Default()

	systemAccount := natsserver.NewAccount("system")
	systemUser := &natsserver.User{
		Username: "system",
		Password: "systempassword",
		Account:  systemAccount,
	}

	regularUser := &natsserver.User{
		Username: "myuser",
		Password: "mypassword",
	}

	opts := &server.Options{
		Options: &natsserver.Options{
			Accounts: []*natsserver.Account{
				systemAccount,
			},
			Users: []*natsserver.User{
				systemUser,
				regularUser,
			},
			SystemAccount: "system",
			JetStream:     true,
			Debug:         debug,
			Trace:         trace,
			StoreDir:      ".nats/jetstream/",
			NoSigs:        true,
			NoLog:         false,
		},
		ReadyTimeout: 5 * time.Second,
	}

	s := server.New(logger, opts)
	err := s.Start()
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.Stop()
}
