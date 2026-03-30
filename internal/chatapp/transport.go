package chatapp

import (
	"context"
	"errors"

	"se-xp-2026-chat/internal/chat"
)

var ErrSessionClosed = errors.New("session is closed")

type Session interface {
	Send(ctx context.Context, msg chat.Message) error
	Recv(ctx context.Context) (chat.Message, error)
	Close() error
}

type Transport interface {
	Listen(ctx context.Context, addr string) (<-chan Session, error)
	Dial(ctx context.Context, peerAddr string) (Session, error)
}
