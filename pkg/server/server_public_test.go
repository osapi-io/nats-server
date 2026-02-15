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
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/suite"

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
	opts := &server.Options{
		Options:      &natsserver.Options{Host: "localhost", Port: 4222},
		ReadyTimeout: 10 * time.Second,
	}
	srv := server.New(s.logger, opts)

	s.NotNil(srv)
	s.Equal(opts, srv.Opts)
}

func (s *ServerPublicTestSuite) TestStart() {
	tests := []struct {
		name        string
		mockSetup   func()
		expectedErr string
	}{
		{
			name: "successfully starts server",
			mockSetup: func() {
				server.NewNATSServer = func(
					_ *natsserver.Options,
				) (server.NATSServerInstance, error) {
					return s.mockNATSServer, nil
				}
				s.mockNATSServer.EXPECT().Start().Times(1)
				s.mockNATSServer.EXPECT().
					ReadyForConnections(gomock.Any()).
					Return(true).
					Times(1)
				s.mockNATSServer.EXPECT().
					SetLogger(gomock.Any(), true, true).
					Times(1)
			},
			expectedErr: "",
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
			expectedErr: "error starting server: invalid options",
		},
		{
			name: "returns error when not ready for connections",
			mockSetup: func() {
				server.NewNATSServer = func(
					_ *natsserver.Options,
				) (server.NATSServerInstance, error) {
					return s.mockNATSServer, nil
				}
				s.mockNATSServer.EXPECT().Start().AnyTimes()
				s.mockNATSServer.EXPECT().
					ReadyForConnections(gomock.Any()).
					Return(false).
					Times(1)
			},
			expectedErr: "server not ready for connections",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			tc.mockSetup()

			err := s.srv.Start()

			if tc.expectedErr == "" {
				s.NoError(err)
			} else {
				s.EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *ServerPublicTestSuite) TestStop() {
	tests := []struct {
		name      string
		setup     func()
		expectErr bool
	}{
		{
			name: "stops running server",
			setup: func() {
				server.NewNATSServer = func(
					_ *natsserver.Options,
				) (server.NATSServerInstance, error) {
					return s.mockNATSServer, nil
				}
				s.mockNATSServer.EXPECT().Start().Times(1)
				s.mockNATSServer.EXPECT().
					ReadyForConnections(gomock.Any()).
					Return(true).
					Times(1)
				s.mockNATSServer.EXPECT().
					SetLogger(gomock.Any(), true, true).
					Times(1)
				s.mockNATSServer.EXPECT().Shutdown().Times(1)
			},
		},
		{
			name:  "handles nil server gracefully",
			setup: func() {},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			originalNewNATSServer := server.NewNATSServer
			defer func() { server.NewNATSServer = originalNewNATSServer }()

			tc.setup()

			if tc.name == "stops running server" {
				err := s.srv.Start()
				s.NoError(err)
			}

			s.NotPanics(func() {
				s.srv.Stop()
			})
		})
	}
}

func TestServerPublicTestSuite(t *testing.T) {
	suite.Run(t, new(ServerPublicTestSuite))
}
