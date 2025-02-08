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
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

// DefaultStoreDir defines the default directory for JetStream storage.
const DefaultStoreDir = "/var/lib/nats/jetstream"

// Option is a function that modifies an Options instance.
type Option func(*Options)

// Options defines the configuration for a NATS server instance.
type Options struct {
	// Enables JetStream.
	JetStream bool
	// Enables tracing.
	Trace bool
	// Enables debug mode.
	Debug bool
	// Host address for the server.
	Host string
	// Port number for the server.
	Port int
	// Disables logging.
	NoLog bool
	// Disables signal handling.
	NoSigs bool
	// Directory path for file-based JetStream storage. If empty, memory storage is used.
	StoreDir string
	// Configurable server readiness timeout
	ReadyTimeout time.Duration
}

// NewOptions initializes an Options struct with defaults and applies functional options.
func NewOptions(opts ...Option) *Options {
	o := &Options{
		JetStream:    true,
		Trace:        false,
		Debug:        false,
		Host:         "localhost",
		Port:         4222,
		NoLog:        false,
		NoSigs:       true,
		StoreDir:     DefaultStoreDir,
		ReadyTimeout: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithJetStream enables or disables JetStream.
func WithJetStream(enabled bool) Option {
	return func(o *Options) {
		o.JetStream = enabled
	}
}

// WithTrace enables or disables trace logging.
func WithTrace(enabled bool) Option {
	return func(o *Options) {
		o.Trace = enabled
	}
}

// WithDebug enables or disables debug mode.
func WithDebug(enabled bool) Option {
	return func(o *Options) {
		o.Debug = enabled
	}
}

// WithHost sets the host.
func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

// WithPort sets the port.
func WithPort(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

// WithNoLog enables or disables logging.
func WithNoLog(enabled bool) Option {
	return func(o *Options) {
		o.NoLog = enabled
	}
}

// WithNoSigs enables or disables signal handling.
func WithNoSigs(enabled bool) Option {
	return func(o *Options) {
		o.NoSigs = enabled
	}
}

// WithStoreDir sets the storage directory.
func WithStoreDir(dir string) Option {
	return func(o *Options) {
		o.StoreDir = dir
	}
}

// WithReadyTimeout sets the timeout for waiting for the server to become ready.
func WithReadyTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.ReadyTimeout = d
	}
}

// ToNATSOptions converts your Options struct to a NATS natsserver.Options struct.
func (o *Options) ToNATSOptions() *natsserver.Options {
	return &natsserver.Options{
		JetStream: o.JetStream,
		Trace:     o.Trace,
		Debug:     o.Debug,
		Host:      o.Host,
		Port:      o.Port,
		NoLog:     o.NoLog,
		NoSigs:    o.NoSigs,
		StoreDir:  o.StoreDir,
	}
}
