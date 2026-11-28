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

package server_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/osapi-io/nats-server/pkg/server"
	"github.com/osapi-io/nats-server/pkg/server/mocks"
)

type ServerPublicTestSuite struct {
	suite.Suite

	mockCtrl       *gomock.Controller
	mockNATSServer *mocks.MockNATSServerInstance
	logger         *slog.Logger
	srv            *server.Server

	// waiting closes when a Start that is never ready waits.
	waiting <-chan struct{}
	// shutDown is true after the mock server has shut down.
	shutDown atomic.Bool
	// addr and logs are the address and the log of a realServer.
	addr string
	logs *lockedBuffer
}

func (s *ServerPublicTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNATSServer = mocks.NewMockNATSServerInstance(s.mockCtrl)
	s.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	s.srv = server.New(s.logger, &server.Options{
		Options:      &natsserver.Options{},
		ReadyTimeout: 5 * time.Second,
	})
	s.waiting = nil
	s.shutDown.Store(false)
	s.addr = ""
	s.logs = nil
}

func (s *ServerPublicTestSuite) SetupSubTest() {
	s.SetupTest()
}

func (s *ServerPublicTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *ServerPublicTestSuite) TestNew() {
	tests := []struct {
		name         string
		opts         *server.Options
		validateFunc func(*server.Server, *server.Options)
	}{
		{
			name: "keeps the options it was given",
			opts: &server.Options{
				Options:      &natsserver.Options{Host: "localhost", Port: 4222},
				ReadyTimeout: 10 * time.Second,
			},
			validateFunc: func(srv *server.Server, opts *server.Options) {
				s.NotNil(srv)
				s.Equal(opts, srv.Opts)
			},
		},
		{
			name: "accepts options with a zero-value embedded config",
			opts: &server.Options{
				Options: &natsserver.Options{},
			},
			validateFunc: func(srv *server.Server, opts *server.Options) {
				s.NotNil(srv)
				s.Equal(opts, srv.Opts)
				s.Zero(srv.Opts.ReadyTimeout)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(server.New(s.logger, tc.opts), tc.opts)
		})
	}
}

func (s *ServerPublicTestSuite) useMock() {
	server.NewNATSServer = func(
		_ *natsserver.Options,
	) (server.NATSServerInstance, error) {
		return s.mockNATSServer, nil
	}
}

func (s *ServerPublicTestSuite) expectReady() {
	s.mockNATSServer.EXPECT().Start().AnyTimes()
	s.mockNATSServer.EXPECT().
		ReadyForConnections(gomock.Any()).
		Return(true).
		Times(1)
	s.mockNATSServer.EXPECT().
		SetLogger(gomock.Any(), false, false).
		Times(1)
}

// expectNotReady makes the mock never ready. Each try waits its timeout,
// and the first one calls logs, if set, with the logger that Start gave
// the server. The channel closes when Start waits.
func (s *ServerPublicTestSuite) expectNotReady(
	logs func(natsserver.Logger),
) <-chan struct{} {
	var (
		logger natsserver.Logger
		once   sync.Once
	)
	waiting := make(chan struct{})

	s.mockNATSServer.EXPECT().
		SetLogger(gomock.Any(), false, false).
		Do(func(l natsserver.Logger, _, _ bool) { logger = l }).
		Times(1)
	s.mockNATSServer.EXPECT().Start().AnyTimes()
	s.mockNATSServer.EXPECT().
		ReadyForConnections(gomock.Any()).
		DoAndReturn(func(d time.Duration) bool {
			once.Do(func() {
				if logs != nil {
					logs(logger)
				}
				close(waiting)
			})
			time.Sleep(d)

			return false
		}).
		AnyTimes()

	return waiting
}

// expectShutdown expects one Shutdown, then one WaitForShutdown.
func (s *ServerPublicTestSuite) expectShutdown() {
	gomock.InOrder(
		s.mockNATSServer.EXPECT().Shutdown().Times(1),
		s.mockNATSServer.EXPECT().WaitForShutdown().Times(1),
	)
}

func (s *ServerPublicTestSuite) TestStart() {
	tests := []struct {
		name string
		// srv is the Server under test. When nil, the test uses s.srv.
		srv       func() *server.Server
		mockSetup func()
		// run starts the server. When nil, the test calls Start with a
		// background context.
		run          func(srv *server.Server) error
		validateFunc func(srv *server.Server, err error)
	}{
		{
			name: "successfully starts server",
			mockSetup: func() {
				s.useMock()
				s.expectReady()
			},
			validateFunc: func(_ *server.Server, err error) {
				s.NoError(err)
			},
		},
		{
			name: "waits the default timeout when ReadyTimeout is zero",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = 0
				s.useMock()
				// With a poll longer than any timeout, the first try waits the
				// whole timeout, which shows the default.
				s.T().Cleanup(server.SetReadyPoll(time.Hour))
				s.mockNATSServer.EXPECT().Start().AnyTimes()
				s.mockNATSServer.EXPECT().
					ReadyForConnections(gomock.Any()).
					DoAndReturn(func(d time.Duration) bool {
						s.InDelta(server.DefaultReadyTimeout, d, float64(time.Second))

						return true
					}).
					Times(1)
				s.mockNATSServer.EXPECT().
					SetLogger(gomock.Any(), false, false).
					Times(1)
			},
			validateFunc: func(_ *server.Server, err error) {
				s.NoError(err)
			},
		},
		{
			name: "refuses a negative ReadyTimeout",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = -time.Second
			},
			validateFunc: func(_ *server.Server, err error) {
				s.ErrorIs(err, server.ErrNegativeReadyTimeout)
				s.EqualError(err, "ready timeout is negative: -1s")
			},
		},
		{
			name: "refuses nil options",
			srv: func() *server.Server {
				return server.New(s.logger, nil)
			},
			mockSetup: func() {},
			validateFunc: func(_ *server.Server, err error) {
				s.ErrorIs(err, server.ErrNoOptions)
			},
		},
		{
			name: "refuses nil NATS options",
			srv: func() *server.Server {
				return server.New(s.logger, &server.Options{})
			},
			mockSetup: func() {},
			validateFunc: func(_ *server.Server, err error) {
				s.ErrorIs(err, server.ErrNoOptions)
			},
		},
		{
			name: "refuses a second Start before Stop",
			mockSetup: func() {
				s.useMock()
				s.expectReady()
			},
			run: func(srv *server.Server) error {
				s.Require().NoError(srv.Start(context.Background()))

				return srv.Start(context.Background())
			},
			validateFunc: func(_ *server.Server, err error) {
				s.ErrorIs(err, server.ErrAlreadyStarted)
			},
		},
		{
			name: "returns error when NewServer fails",
			mockSetup: func() {
				server.NewNATSServer = func(
					_ *natsserver.Options,
				) (server.NATSServerInstance, error) {
					return nil, errors.New("invalid options")
				}
			},
			validateFunc: func(_ *server.Server, err error) {
				s.EqualError(err, "error starting server: invalid options")
			},
		},
		{
			name: "attaches the logger before the server runs, with Debug and Trace",
			srv: func() *server.Server {
				return server.New(s.logger, &server.Options{
					Options: &natsserver.Options{Debug: true, Trace: true},
				})
			},
			mockSetup: func() {
				s.useMock()
				gomock.InOrder(
					s.mockNATSServer.EXPECT().
						SetLogger(gomock.Any(), true, true).
						Times(1),
					s.mockNATSServer.EXPECT().Start().Times(1),
				)
				s.mockNATSServer.EXPECT().
					ReadyForConnections(gomock.Any()).
					DoAndReturn(func(time.Duration) bool {
						time.Sleep(50 * time.Millisecond)

						return true
					}).
					Times(1)
				s.expectShutdown()
			},
			validateFunc: func(srv *server.Server, err error) {
				s.NoError(err)
				srv.Stop()
			},
		},
		{
			name: "a nil logger is slog.Default",
			srv: func() *server.Server {
				return server.New(nil, &server.Options{
					Options: &natsserver.Options{},
				})
			},
			mockSetup: func() {
				s.useMock()
				s.expectReady()
				s.expectShutdown()
			},
			validateFunc: func(srv *server.Server, err error) {
				s.NoError(err)
				srv.Stop()
			},
		},
		{
			name: "a Server that New did not make logs to slog.Default",
			srv: func() *server.Server {
				return &server.Server{Opts: &server.Options{Options: &natsserver.Options{}}}
			},
			mockSetup: func() {
				s.useMock()
				s.expectReady()
				s.expectShutdown()
			},
			validateFunc: func(srv *server.Server, err error) {
				s.NoError(err)
				s.NotPanics(srv.Stop)
			},
		},
		{
			name: "shuts the server down when it is not ready",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = 100 * time.Millisecond
				s.useMock()
				s.expectNotReady(nil)
				s.expectShutdown()
			},
			validateFunc: func(_ *server.Server, err error) {
				s.EqualError(err, "server not ready for connections after 100ms")
			},
		},
		{
			name: "returns the fatal error that the server logged at once",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = time.Hour
				s.useMock()
				s.expectNotReady(func(l natsserver.Logger) {
					l.Fatalf("Error listening on port: %s", "127.0.0.1:4222")
				})
				s.expectShutdown()
			},
			validateFunc: func(_ *server.Server, err error) {
				s.EqualError(err, "server did not start: Error listening on port: 127.0.0.1:4222")
			},
		},
		{
			name: "keeps a % in a fatal message that has no arguments",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = time.Hour
				s.useMock()
				s.expectNotReady(func(l natsserver.Logger) {
					// A func value, so that vet takes the text as a message
					// and not as a format, as the NATS server's own calls do.
					fatalf := l.Fatalf
					fatalf("store is 100% full")
				})
				s.expectShutdown()
			},
			validateFunc: func(_ *server.Server, err error) {
				s.EqualError(err, "server did not start: store is 100% full")
			},
		},
		{
			name: "names the last error when the server is not ready",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = 100 * time.Millisecond
				s.useMock()
				s.expectNotReady(func(l natsserver.Logger) {
					l.Errorf("first")
					l.Errorf("the %s one", "last")
				})
				s.expectShutdown()
			},
			validateFunc: func(_ *server.Server, err error) {
				s.EqualError(
					err,
					"server not ready for connections after 100ms: the last one",
				)
			},
		},
		{
			name: "keeps a % in an error message that has no arguments",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = 100 * time.Millisecond
				s.useMock()
				s.expectNotReady(func(l natsserver.Logger) {
					// A func value, so that vet takes the text as a message
					// and not as a format, as the NATS server's own calls do.
					errorf := l.Errorf
					errorf("disk 100% full")
				})
				s.expectShutdown()
			},
			validateFunc: func(_ *server.Server, err error) {
				s.EqualError(
					err,
					"server not ready for connections after 100ms: disk 100% full",
				)
			},
		},
		{
			name: "a ctx that ends stops a Start that waits",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = time.Hour
				s.useMock()
				s.expectNotReady(nil)
				s.expectShutdown()
			},
			run: func(srv *server.Server) error {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				return srv.Start(ctx)
			},
			validateFunc: func(_ *server.Server, err error) {
				s.EqualError(err, "server start stopped: context canceled")
			},
		},
		{
			name: "a ctx that ends as the server becomes ready shuts it down",
			mockSetup: func() {
				s.useMock()
				s.mockNATSServer.EXPECT().
					SetLogger(gomock.Any(), false, false).
					Times(1)
				s.mockNATSServer.EXPECT().Start().AnyTimes()
				s.expectShutdown()
			},
			run: func(srv *server.Server) error {
				ctx, cancel := context.WithCancel(context.Background())
				s.mockNATSServer.EXPECT().
					ReadyForConnections(gomock.Any()).
					DoAndReturn(func(time.Duration) bool {
						cancel()

						return true
					}).
					Times(1)

				return srv.Start(ctx)
			},
			validateFunc: func(srv *server.Server, err error) {
				s.EqualError(err, "server start stopped: context canceled")
				s.Empty(srv.ClientURL())
			},
		},
		{
			name: "Stop ends a Start that waits, and a second Start is refused meanwhile",
			mockSetup: func() {
				s.srv.Opts.ReadyTimeout = time.Hour
				s.useMock()
				s.waiting = s.expectNotReady(nil)
				s.expectShutdown()
			},
			run: func(srv *server.Server) error {
				done := make(chan error, 1)
				go func() { done <- srv.Start(context.Background()) }()

				<-s.waiting
				s.ErrorIs(srv.Start(context.Background()), server.ErrAlreadyStarted)
				srv.Stop()

				return <-done
			},
			validateFunc: func(_ *server.Server, err error) {
				s.EqualError(err, "server start stopped: context canceled")
			},
		},
		{
			name: "a real server starts, accepts a connection, and frees its port on Stop",
			srv: func() *server.Server {
				return s.realServer(s.freePort(), 0)
			},
			mockSetup: func() {},
			validateFunc: func(srv *server.Server, err error) {
				s.Require().NoError(err)

				conn, err := net.Dial("tcp", s.addr)
				s.Require().NoError(err)
				s.Require().NoError(conn.Close())

				srv.Stop()

				l, err := net.Listen("tcp", s.addr)
				s.Require().NoError(err)
				s.NoError(l.Close())
				s.Contains(s.logs.String(), "Starting nats-server")
			},
		},
		{
			name: "a real server with port -1 has the port that it listens on in ClientURL",
			srv: func() *server.Server {
				return s.realServer(-1, 0)
			},
			mockSetup: func() {},
			validateFunc: func(srv *server.Server, err error) {
				s.Require().NoError(err)
				defer srv.Stop()

				url := srv.ClientURL()
				s.Require().NotEmpty(url)
				s.NotContains(url, ":-1")

				conn, err := net.Dial("tcp", url[len("nats://"):])
				s.Require().NoError(err)
				s.NoError(conn.Close())
			},
		},
		{
			name: "a real server on a port in use returns the cause at once, not after ReadyTimeout",
			srv: func() *server.Server {
				held, err := net.Listen("tcp", "127.0.0.1:0")
				s.Require().NoError(err)
				s.T().Cleanup(func() { s.NoError(held.Close()) })

				return s.realServer(held.Addr().(*net.TCPAddr).Port, time.Minute)
			},
			mockSetup: func() {},
			run: func(srv *server.Server) error {
				begin := time.Now()
				err := srv.Start(context.Background())
				s.Less(time.Since(begin), 10*time.Second)

				return err
			},
			validateFunc: func(srv *server.Server, err error) {
				s.Require().Error(err)
				s.Contains(err.Error(), "server did not start:")
				s.Contains(err.Error(), "listening")
				s.Empty(srv.ClientURL())
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			srv := s.srv
			if tc.srv != nil {
				srv = tc.srv()
			}
			tc.mockSetup()

			run := tc.run
			if run == nil {
				run = func(srv *server.Server) error {
					return srv.Start(context.Background())
				}
			}

			tc.validateFunc(srv, run(srv))
		})
	}
}

func (s *ServerPublicTestSuite) TestStop() {
	tests := []struct {
		name         string
		setup        func()
		validateFunc func()
	}{
		{
			name: "shuts down, waits, and lets Start run again",
			setup: func() {
				s.useMock()
				s.mockNATSServer.EXPECT().Start().AnyTimes()
				s.mockNATSServer.EXPECT().
					ReadyForConnections(gomock.Any()).
					Return(true).
					Times(2)
				s.mockNATSServer.EXPECT().
					SetLogger(gomock.Any(), false, false).
					Times(2)
				s.expectShutdown()

				s.Require().NoError(s.srv.Start(context.Background()))
			},
			validateFunc: func() {
				s.srv.Stop()
				s.srv.Stop()
				s.NoError(s.srv.Start(context.Background()))
			},
		},
		{
			name:  "handles nil server gracefully",
			setup: func() {},
			validateFunc: func() {
				s.NotPanics(s.srv.Stop)
			},
		},
		{
			name: "returns after the Start in progress has shut its server down",
			setup: func() {
				s.srv.Opts.ReadyTimeout = time.Hour
				s.useMock()
				s.waiting = s.expectNotReady(nil)
				gomock.InOrder(
					s.mockNATSServer.EXPECT().Shutdown().Times(1),
					s.mockNATSServer.EXPECT().
						WaitForShutdown().
						Do(func() { s.shutDown.Store(true) }).
						Times(1),
				)
			},
			validateFunc: func() {
				done := make(chan error, 1)
				go func() { done <- s.srv.Start(context.Background()) }()

				<-s.waiting
				s.srv.Stop()

				s.True(s.shutDown.Load(), "the server still runs after Stop")
				select {
				case err := <-done:
					s.EqualError(err, "server start stopped: context canceled")
				default:
					s.Fail("Start still runs after Stop")
				}
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			tc.setup()

			tc.validateFunc()
		})
	}
}

func (s *ServerPublicTestSuite) TestClientURL() {
	tests := []struct {
		name         string
		start        bool
		validateFunc func(string)
	}{
		{
			name:  "is the running server's URL",
			start: true,
			validateFunc: func(url string) {
				s.Equal("nats://127.0.0.1:4222", url)
			},
		},
		{
			name: "is empty when the server does not run",
			validateFunc: func(url string) {
				s.Empty(url)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			if tc.start {
				s.useMock()
				s.expectReady()
				s.mockNATSServer.EXPECT().ClientURL().Return("nats://127.0.0.1:4222")
				s.Require().NoError(s.srv.Start(context.Background()))
			}

			tc.validateFunc(s.srv.ClientURL())
		})
	}
}

// freePort is a port that nothing listens on now.
func (s *ServerPublicTestSuite) freePort() int {
	free, err := net.Listen("tcp", "127.0.0.1:0")
	s.Require().NoError(err)
	port := free.Addr().(*net.TCPAddr).Port
	s.Require().NoError(free.Close())

	return port
}

// realServer is a Server that runs a real NATS server on port. It sets
// s.addr to the address of port and s.logs to what the server logs.
func (s *ServerPublicTestSuite) realServer(
	port int,
	readyTimeout time.Duration,
) *server.Server {
	s.addr = net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	s.logs = &lockedBuffer{buf: &bytes.Buffer{}}

	return server.New(
		slog.New(slog.NewTextHandler(s.logs, nil)),
		&server.Options{
			Options: &natsserver.Options{
				Host:   "127.0.0.1",
				Port:   port,
				NoSigs: true,
			},
			ReadyTimeout: readyTimeout,
		},
	)
}

// lockedBuffer is a bytes.Buffer that the server's goroutines can write to
// at the same time.
type lockedBuffer struct {
	mu  sync.Mutex
	buf *bytes.Buffer
}

func (b *lockedBuffer) Write(
	p []byte,
) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.Write(p)
}

// String is what the server wrote so far.
func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.String()
}

func TestServerPublicTestSuite(t *testing.T) {
	suite.Run(t, new(ServerPublicTestSuite))
}
