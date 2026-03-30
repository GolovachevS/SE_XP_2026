package chatapp

import (
	"context"

	"se-xp-2026-chat/internal/chat"
)

// NoopTransport is a placeholder transport kept for early-stage tests and experiments.
type NoopTransport struct{}

// NewNoopTransport creates a transport that never performs real network I/O.
func NewNoopTransport() NoopTransport {
	return NoopTransport{}
}

func (NoopTransport) Listen(context.Context, string) (<-chan Session, error) {
	return make(chan Session), nil
}

func (NoopTransport) Dial(context.Context, string) (Session, error) {
	return noopSession{}, nil
}

type noopSession struct{}

func (noopSession) Send(context.Context, chat.Message) error {
	return nil
}

func (noopSession) Recv(context.Context) (chat.Message, error) {
	return chat.Message{}, ErrSessionClosed
}

func (noopSession) Close() error {
	return nil
}
