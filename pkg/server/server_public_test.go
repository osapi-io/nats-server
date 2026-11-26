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
	"errors"
	"log/slog"
	"net"
	"os"
	"strconv"
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
		SetLogger(gomock.Any(), true, true).
		Times(1)
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
				s.mockNATSServer.EXPECT().Start().AnyTimes()
				s.mockNATSServer.EXPECT().
					ReadyForConnections(server.DefaultReadyTimeout).
					Return(true).
					Times(1)
				s.mockNATSServer.EXPECT().
					SetLogger(gomock.Any(), true, true).
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
			readyTimeout: 5 * time.Second,
			mockSetup: func() {
				s.useMock()
				s.mockNATSServer.EXPECT().Start().AnyTimes()
				s.mockNATSServer.EXPECT().
					ReadyForConnections(gomock.Any()).
					Return(false).
					Times(1)
				gomock.InOrder(
					s.mockNATSServer.EXPECT().Shutdown().Times(1),
					s.mockNATSServer.EXPECT().WaitForShutdown().Times(1),
				)
			},
			validateFunc: func(err error) {
				s.EqualError(err, "server not ready for connections")
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			s.srv.Opts.ReadyTimeout = tc.readyTimeout
			tc.mockSetup()

			tc.validateFunc(s.srv.Start())
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
			s.Require().NoError(s.srv.Start())

			tc.validateFunc(s.srv.Start())
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
					SetLogger(gomock.Any(), true, true).
					Times(2)
				gomock.InOrder(
					s.mockNATSServer.EXPECT().Shutdown().Times(1),
					s.mockNATSServer.EXPECT().WaitForShutdown().Times(1),
				)

				s.Require().NoError(s.srv.Start())
			},
			validateFunc: func() {
				s.srv.Stop()
				s.srv.Stop()
				s.NoError(s.srv.Start())
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

func (s *ServerPublicTestSuite) TestRealServer() {
	tests := []struct {
		name         string
		validateFunc func(*server.Server, string)
	}{
		{
			name: "starts, accepts a connection, and frees its port on Stop",
			validateFunc: func(srv *server.Server, addr string) {
				s.Require().NoError(srv.Start())

				conn, err := net.Dial("tcp", addr)
				s.Require().NoError(err)
				s.Require().NoError(conn.Close())

				srv.Stop()

				l, err := net.Listen("tcp", addr)
				s.Require().NoError(err)
				s.NoError(l.Close())
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			free, err := net.Listen("tcp", "127.0.0.1:0")
			s.Require().NoError(err)
			port := free.Addr().(*net.TCPAddr).Port
			s.Require().NoError(free.Close())

			srv := server.New(s.logger, &server.Options{
				Options: &natsserver.Options{
					Host:   "127.0.0.1",
					Port:   port,
					NoSigs: true,
				},
			})

			tc.validateFunc(srv, net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
		})
	}
}

func TestServerPublicTestSuite(t *testing.T) {
	suite.Run(t, new(ServerPublicTestSuite))
}
