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

	"github.com/nats-io/nats.go"
)

// StreamOption is a function that modifies a StreamOptions instance.
type StreamOption func(*StreamOptions)

// StreamOptions defines the configuration for a NATS JetStream stream.
type StreamOptions struct {
	// Stream name.
	Name string
	// List of subjects associated with the stream.
	Subjects []string
	// Storage type (e.g., MemoryStorage or FileStorage).
	Storage nats.StorageType
	// Number of data replicas.
	Replicas int
	// Consumers tied to this stream.
	Consumers []*ConsumerOptions
}

// ConsumerOption is a function that modifies a ConsumerOptions instance.
type ConsumerOption func(*ConsumerOptions)

// ConsumerOptions defines the configuration for a NATS JetStream consumer.
type ConsumerOptions struct {
	// Durable name for the consumer.
	Durable string
	// Acknowledgment policy.
	AckPolicy nats.AckPolicy
	// Maximum number of times a message is delivered.
	MaxDeliver int
	// Time to wait for an acknowledgment.
	AckWait time.Duration
	// Deliver policy (e.g., All, New, Last, etc.).
	DeliverPolicy nats.DeliverPolicy
	// Maximum number of unacknowledged messages allowed.
	MaxAckPending int
}

// NewStreamOptions initializes a StreamOptions struct with defaults and applies functional options.
func NewStreamOptions(opts ...StreamOption) *StreamOptions {
	o := &StreamOptions{
		Storage:  nats.FileStorage,
		Replicas: 1,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithStreamName sets the name of the stream.
func WithStreamName(name string) StreamOption {
	return func(o *StreamOptions) {
		o.Name = name
	}
}

// WithSubjects sets the subjects for the stream.
func WithSubjects(subjects ...string) StreamOption {
	return func(o *StreamOptions) {
		o.Subjects = subjects
	}
}

// WithStorage sets the storage type for the stream (e.g., FileStorage or MemoryStorage).
func WithStorage(storage nats.StorageType) StreamOption {
	return func(o *StreamOptions) {
		o.Storage = storage
	}
}

// WithReplicas sets the number of replicas for the stream.
func WithReplicas(replicas int) StreamOption {
	return func(o *StreamOptions) {
		o.Replicas = replicas
	}
}

// WithConsumer ties a ConsumerOptions instance to a StreamOptions instance.
func WithConsumer(consumer *ConsumerOptions) StreamOption {
	return func(o *StreamOptions) {
		o.Consumers = append(o.Consumers, consumer)
	}
}

// ToNATS converts StreamOptions to a nats.StreamConfig.
func (o *StreamOptions) ToNATS() *nats.StreamConfig {
	return &nats.StreamConfig{
		Name:     o.Name,
		Subjects: o.Subjects,
		Storage:  o.Storage,
		Replicas: o.Replicas,
	}
}

// NewConsumerOptions initializes a ConsumerOptions struct with defaults and applies functional options.
func NewConsumerOptions(opts ...ConsumerOption) *ConsumerOptions {
	o := &ConsumerOptions{
		AckPolicy:  nats.AckExplicitPolicy,
		MaxDeliver: 5,
		AckWait:    30 * time.Second,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithDurable sets the durable name of the consumer.
func WithDurable(name string) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.Durable = name
	}
}

// WithAckPolicy sets the acknowledgment policy of the consumer.
func WithAckPolicy(policy nats.AckPolicy) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.AckPolicy = policy
	}
}

// WithMaxDeliver sets the max number of delivery attempts for a message.
func WithMaxDeliver(max int) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.MaxDeliver = max
	}
}

// WithAckWait sets the acknowledgment wait duration.
func WithAckWait(d time.Duration) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.AckWait = d
	}
}

// WithDeliverPolicy sets the deliver policy for a consumer.
func WithDeliverPolicy(policy nats.DeliverPolicy) ConsumerOption {
	return func(opts *ConsumerOptions) {
		opts.DeliverPolicy = policy
	}
}

// WithMaxAckPending sets the maximum number of unacknowledged messages a consumer can have.
func WithMaxAckPending(max int) ConsumerOption {
	return func(opts *ConsumerOptions) {
		opts.MaxAckPending = max
	}
}

// ToNATS converts ConsumerOptions to a nats.ConsumerConfig.
func (o *ConsumerOptions) ToNATS() *nats.ConsumerConfig {
	return &nats.ConsumerConfig{
		Durable:       o.Durable,
		AckPolicy:     o.AckPolicy,
		AckWait:       o.AckWait,
		MaxDeliver:    o.MaxDeliver,
		DeliverPolicy: o.DeliverPolicy,
		MaxAckPending: o.MaxAckPending,
	}
}
