package events

import (
	"context"
	"time"
)

type Metadata struct {
	At       time.Time `json:"at"`
	UserId   string    `json:"user_id"`
	ObjectId string    `json:"object_id"`
}

type Message struct {
	Topic    Topic     `json:"topic"`
	Metadata *Metadata `json:"metadata"`
}

type MessageParams struct {
	UserId   string
	ObjectId string
	Topic    Topic
}

type EventQueue interface {
	Push(ctx context.Context, message *Message) error
	Close() error
}

type option func(*Message)

func WithTopic(t Topic) option {
	return func(m *Message) {
		m.Topic = t
	}
}

func WithUser(id string) option {
	return func(m *Message) {
		m.Metadata.UserId = id
	}
}

func WithTarget(id string) option {
	return func(m *Message) {
		m.Metadata.ObjectId = id
	}
}

func NewMessage(opts ...option) *Message {
	m := &Message{
		Metadata: &Metadata{},
	}
	for _, opt := range opts {
		opt(m)
	}

	return m
}
