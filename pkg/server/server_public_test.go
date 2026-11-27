// Copyright (c) 2026 John Dewey

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
}

func (s *ServerPublicTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNATSServer = mocks.NewMockNATSServerInstance(s.mockCtrl)
	s.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	s.srv = server.New(s.logger, &server.Options{
		Options:      &natsserver.Options{},
		ReadyTimeout: 5 * time.Second,
	})
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
		name         string
		readyTimeout time.Duration
		mockSetup    func()
		validateFunc func(error)
	}{
		{
			name:         "successfully starts server",
			readyTimeout: 5 * time.Second,
			mockSetup: func() {
				s.useMock()
				s.expectReady()
			},
			validateFunc: func(err error) {
				s.NoError(err)
			},
		},
		{
			name: "waits the default timeout when ReadyTimeout is zero",
			mockSetup: func() {
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
			validateFunc: func(err error) {
				s.NoError(err)
			},
		},
		{
			name:         "refuses a negative ReadyTimeout",
			readyTimeout: -time.Second,
			mockSetup:    func() {},
			validateFunc: func(err error) {
				s.EqualError(err, "ready timeout -1s is negative")
			},
		},
		{
			name:         "returns error when NewServer fails",
			readyTimeout: 5 * time.Second,
			mockSetup: func() {
				server.NewNATSServer = func(
					_ *natsserver.Options,
				) (server.NATSServerInstance, error) {
					return nil, errors.New("invalid options")
				}
			},
			validateFunc: func(err error) {
				s.EqualError(err, "error starting server: invalid options")
			},
		},
		{
			name:         "shuts the server down when it is not ready",
			readyTimeout: 100 * time.Millisecond,
			mockSetup: func() {
				s.useMock()
				s.expectNotReady(nil)
				s.expectShutdown()
			},
			validateFunc: func(err error) {
				s.EqualError(err, "server not ready for connections after 100ms")
			},
		},
		{
			name:         "returns the fatal error that the server logged at once",
			readyTimeout: time.Hour,
			mockSetup: func() {
				s.useMock()
				s.expectNotReady(func(l natsserver.Logger) {
					l.Fatalf("Error listening on port: %s", "127.0.0.1:4222")
				})
				s.expectShutdown()
			},
			validateFunc: func(err error) {
				s.EqualError(err, "server did not start: Error listening on port: 127.0.0.1:4222")
			},
		},
		{
			name:         "names the last error when the server is not ready",
			readyTimeout: 100 * time.Millisecond,
			mockSetup: func() {
				s.useMock()
				s.expectNotReady(func(l natsserver.Logger) {
					l.Errorf("first")
					l.Errorf("the %s one", "last")
				})
				s.expectShutdown()
			},
			validateFunc: func(err error) {
				s.EqualError(
					err,
					"server not ready for connections after 100ms: the last one",
				)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			s.srv.Opts.ReadyTimeout = tc.readyTimeout
			tc.mockSetup()

			tc.validateFunc(s.srv.Start(context.Background()))
		})
	}
}

func (s *ServerPublicTestSuite) TestStartOptions() {
	tests := []struct {
		name         string
		srv          func() *server.Server
		mockSetup    func()
		validateFunc func(error)
	}{
		{
			name: "refuses nil options",
			srv: func() *server.Server {
				return server.New(s.logger, nil)
			},
			mockSetup: func() {},
			validateFunc: func(err error) {
				s.EqualError(err, "server options are nil")
			},
		},
		{
			name: "refuses nil NATS options",
			srv: func() *server.Server {
				return server.New(s.logger, &server.Options{})
			},
			mockSetup: func() {},
			validateFunc: func(err error) {
				s.EqualError(err, "server options are nil")
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
				s.mockNATSServer.EXPECT().Shutdown().Times(1)
				s.mockNATSServer.EXPECT().WaitForShutdown().Times(1)
			},
			validateFunc: func(err error) {
				s.NoError(err)
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
				s.mockNATSServer.EXPECT().Shutdown().Times(1)
				s.mockNATSServer.EXPECT().WaitForShutdown().Times(1)
			},
			validateFunc: func(err error) {
				s.NoError(err)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			tc.mockSetup()
			srv := tc.srv()

			tc.validateFunc(srv.Start(context.Background()))
			srv.Stop()
		})
	}
}

func (s *ServerPublicTestSuite) TestStartTwice() {
	tests := []struct {
		name         string
		validateFunc func(error)
	}{
		{
			name: "refuses a second Start before Stop",
			validateFunc: func(err error) {
				s.EqualError(err, "server already started")
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			s.useMock()
			s.expectReady()
			s.Require().NoError(s.srv.Start(context.Background()))

			tc.validateFunc(s.srv.Start(context.Background()))
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
				gomock.InOrder(
					s.mockNATSServer.EXPECT().Shutdown().Times(1),
					s.mockNATSServer.EXPECT().WaitForShutdown().Times(1),
				)

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

func (s *ServerPublicTestSuite) TestStartStops() {
	tests := []struct {
		name         string
		run          func(waiting <-chan struct{}) error
		validateFunc func(error)
	}{
		{
			name: "a ctx that ends stops a Start that waits",
			run: func(<-chan struct{}) error {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				return s.srv.Start(ctx)
			},
			validateFunc: func(err error) {
				s.EqualError(err, "server start stopped: context canceled")
			},
		},
		{
			name: "Stop ends a Start that waits, and a second Start is refused meanwhile",
			run: func(waiting <-chan struct{}) error {
				done := make(chan error, 1)
				go func() { done <- s.srv.Start(context.Background()) }()

				<-waiting
				s.EqualError(s.srv.Start(context.Background()), "server already started")
				s.srv.Stop()

				return <-done
			},
			validateFunc: func(err error) {
				s.EqualError(err, "server start stopped: context canceled")
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			s.srv.Opts.ReadyTimeout = time.Hour
			s.useMock()
			waiting := s.expectNotReady(nil)
			s.expectShutdown()

			tc.validateFunc(tc.run(waiting))
		})
	}
}

func (s *ServerPublicTestSuite) TestStopAsStartBecomesReady() {
	tests := []struct {
		name         string
		validateFunc func(err error)
	}{
		{
			name: "a Stop while the server becomes ready shuts it down",
			validateFunc: func(err error) {
				s.EqualError(err, "server start stopped: context canceled")
				s.Empty(s.srv.ClientURL())
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			entered := make(chan struct{})
			stopped := make(chan struct{})

			s.useMock()
			s.mockNATSServer.EXPECT().
				SetLogger(gomock.Any(), false, false).
				Times(1)
			s.mockNATSServer.EXPECT().Start().AnyTimes()
			s.mockNATSServer.EXPECT().
				ReadyForConnections(gomock.Any()).
				DoAndReturn(func(time.Duration) bool {
					close(entered)
					<-stopped

					return true
				}).
				Times(1)
			s.expectShutdown()

			done := make(chan error, 1)
			go func() { done <- s.srv.Start(context.Background()) }()

			<-entered
			s.srv.Stop()
			close(stopped)

			tc.validateFunc(<-done)
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

func (s *ServerPublicTestSuite) TestZeroServer() {
	tests := []struct {
		name string
	}{
		{name: "a Server that New did not make logs to slog.Default"},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			s.useMock()
			s.expectReady()
			s.expectShutdown()
			srv := &server.Server{Opts: &server.Options{Options: &natsserver.Options{}}}

			s.NotPanics(func() {
				s.NoError(srv.Start(context.Background()))
				srv.Stop()
			})
		})
	}
}

func (s *ServerPublicTestSuite) TestRealServer() {
	tests := []struct {
		name         string
		port         func() (int, func())
		readyTimeout time.Duration
		validateFunc func(*server.Server, string, func() string)
	}{
		{
			name: "starts, accepts a connection, and frees its port on Stop",
			port: s.freePort,
			validateFunc: func(srv *server.Server, addr string, logs func() string) {
				s.Require().NoError(srv.Start(context.Background()))

				conn, err := net.Dial("tcp", addr)
				s.Require().NoError(err)
				s.Require().NoError(conn.Close())

				srv.Stop()

				l, err := net.Listen("tcp", addr)
				s.Require().NoError(err)
				s.NoError(l.Close())
				s.Contains(logs(), "Starting nats-server")
			},
		},
		{
			name: "with port -1, ClientURL is the port that it listens on",
			port: func() (int, func()) { return -1, func() {} },
			validateFunc: func(srv *server.Server, _ string, _ func() string) {
				s.Require().NoError(srv.Start(context.Background()))
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
			name: "a port in use returns the cause at once, not after ReadyTimeout",
			port: func() (int, func()) {
				held, err := net.Listen("tcp", "127.0.0.1:0")
				s.Require().NoError(err)

				return held.Addr().(*net.TCPAddr).Port, func() { s.NoError(held.Close()) }
			},
			readyTimeout: time.Minute,
			validateFunc: func(srv *server.Server, _ string, _ func() string) {
				begin := time.Now()
				err := srv.Start(context.Background())

				s.Require().Error(err)
				s.Contains(err.Error(), "server did not start:")
				s.Contains(err.Error(), "listening")
				s.Less(time.Since(begin), 10*time.Second)
				s.Empty(srv.ClientURL())
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			port, release := tc.port()
			defer release()

			logs := &lockedBuffer{buf: &bytes.Buffer{}}
			srv := server.New(
				slog.New(slog.NewTextHandler(logs, nil)),
				&server.Options{
					Options: &natsserver.Options{
						Host:   "127.0.0.1",
						Port:   port,
						NoSigs: true,
					},
					ReadyTimeout: tc.readyTimeout,
				},
			)

			addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
			tc.validateFunc(srv, addr, logs.String)
		})
	}
}

// freePort is a port that nothing listens on now.
func (s *ServerPublicTestSuite) freePort() (int, func()) {
	free, err := net.Listen("tcp", "127.0.0.1:0")
	s.Require().NoError(err)
	port := free.Addr().(*net.TCPAddr).Port
	s.Require().NoError(free.Close())

	return port, func() {}
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
