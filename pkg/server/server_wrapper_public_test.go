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
	"testing"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/nats-server/pkg/server"
)

type ServerWrapperPublicTestSuite struct {
	suite.Suite
}

func (s *ServerWrapperPublicTestSuite) TestNewNATSServer() {
	tests := []struct {
		name         string
		opts         *natsserver.Options
		validateFunc func(server.NATSServerInstance, error)
	}{
		{
			name: "makes a server from valid options",
			opts: &natsserver.Options{Host: "127.0.0.1", Port: -1, NoSigs: true},
			validateFunc: func(srv server.NATSServerInstance, err error) {
				s.Require().NoError(err)
				s.NotNil(srv)
			},
		},
		{
			name: "returns the error of options that the NATS server refuses",
			opts: &natsserver.Options{
				NoSigs:    true,
				Websocket: natsserver.WebsocketOpts{Port: 8080},
			},
			validateFunc: func(srv server.NATSServerInstance, err error) {
				s.EqualError(err, "websocket requires TLS configuration")
				s.Nil(srv)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(server.NewNATSServer(tc.opts))
		})
	}
}

func TestServerWrapperPublicTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWrapperPublicTestSuite))
}
