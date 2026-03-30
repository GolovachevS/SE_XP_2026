package app

import (
	"context"
	"errors"

	"se-xp-2026-chat/internal/chat"
)

// ErrSessionClosed signals that the underlying transport stream was closed.
var ErrSessionClosed = errors.New("session is closed")

// Session represents one active chat stream between two peers.
type Session interface {
	// Send writes a single chat message to the remote peer.
	Send(ctx context.Context, msg chat.Message) error
	// Recv blocks until a message is received or the session finishes.
	Recv(ctx context.Context) (chat.Message, error)
	// Close releases transport resources associated with the stream.
	Close() error
}

// Transport hides the concrete networking implementation used by the app.
type Transport interface {
	// Listen starts accepting incoming chat sessions on addr.
	Listen(ctx context.Context, addr string) (<-chan Session, error)
	// Dial establishes an outgoing chat session to peerAddr.
	Dial(ctx context.Context, peerAddr string) (Session, error)
}
