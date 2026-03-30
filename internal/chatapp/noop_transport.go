package chatapp

import (
	"context"

	"se-xp-2026-chat/internal/chat"
)

type NoopTransport struct{}

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
