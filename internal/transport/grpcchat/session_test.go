package grpcchat

import (
	"context"
	"errors"
	"testing"
	"time"

	apichatv1 "se-xp-2026-chat/api/chat/v1"
	"se-xp-2026-chat/internal/chat"
	"se-xp-2026-chat/internal/chatapp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stubStream struct {
	sendErr error
	recvErr error
}

func (s stubStream) Send(*apichatv1.Envelope) error {
	return s.sendErr
}

func (s stubStream) Recv() (*apichatv1.Envelope, error) {
	return nil, s.recvErr
}

func (s stubStream) Context() context.Context {
	return context.Background()
}

func TestSessionRecv_TreatsUnavailableAsSessionClosed(t *testing.T) {
	t.Parallel()

	s := &Session{stream: stubStream{recvErr: status.Error(codes.Unavailable, "server down")}, closeFn: func() error { return nil }}

	_, err := s.Recv(context.Background())
	if !errors.Is(err, chatapp.ErrSessionClosed) {
		t.Fatalf("expected ErrSessionClosed, got %v", err)
	}
}

func TestSessionSend_TreatsUnavailableAsSessionClosed(t *testing.T) {
	t.Parallel()

	s := &Session{stream: stubStream{sendErr: status.Error(codes.Unavailable, "server down")}, closeFn: func() error { return nil }}

	msg := chat.Message{Sender: "Alice", SentAt: time.Now(), Text: "hi"}
	err := s.Send(context.Background(), msg)
	if !errors.Is(err, chatapp.ErrSessionClosed) {
		t.Fatalf("expected ErrSessionClosed, got %v", err)
	}
}
