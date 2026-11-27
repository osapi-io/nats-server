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

// Package main is an example that starts an embedded NATS server that
// authenticates clients with NKeys.
package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nkeys"
	"github.com/osapi-io/nats-server/pkg/server"
)

// publicKey reads the seed that `nk -gen user -pubout` wrote on the first line
// of path, and returns its public key. The README makes the files.
func publicKey(
	path string,
) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s, make it as the README says: %w", path, err)
	}

	line, _, _ := bufio.NewReader(bytes.NewReader(data)).ReadLine()
	kp, err := nkeys.FromSeed(line)
	if err != nil {
		return "", fmt.Errorf("parse the seed in %s: %w", path, err)
	}

	return kp.PublicKey()
}

func main() {
	debug := true
	trace := debug
	logger := slog.Default()

	service1, err := publicKey(".nkeys/service1.nk")
	if err != nil {
		logger.Error("failed to read the key of service 1", "error", err)
		os.Exit(1)
	}

	service2, err := publicKey(".nkeys/service2.nk")
	if err != nil {
		logger.Error("failed to read the key of service 2", "error", err)
		os.Exit(1)
	}

	nkeyUsers := []*natsserver.NkeyUser{
		{
			Nkey: service1,
		},
		{
			Nkey: service2, // disabled
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
	err = s.Start(context.Background())
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.Stop()
}
