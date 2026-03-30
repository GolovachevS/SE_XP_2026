package grpcchat

import (
	"context"
	"errors"
	"io"

	apichatv1 "se-xp-2026-chat/api/chat/v1"
	chatapp "se-xp-2026-chat/internal/app"
	"se-xp-2026-chat/internal/chat"
	protocolchatv1 "se-xp-2026-chat/internal/protocol/chatv1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type messageStream interface {
	Send(*apichatv1.Envelope) error
	Recv() (*apichatv1.Envelope, error)
	Context() context.Context
}

type Session struct {
	stream  messageStream
	closeFn func() error
}

func NewClientSession(stream apichatv1.ChatService_ChatStreamClient) *Session {
	return &Session{
		stream:  stream,
		closeFn: stream.CloseSend,
	}
}

func NewServerSession(stream apichatv1.ChatService_ChatStreamServer) *Session {
	return &Session{
		stream:  stream,
		closeFn: func() error { return nil },
	}
}

func (s *Session) Send(ctx context.Context, msg chat.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := s.stream.Send(protocolchatv1.ToEnvelope(msg)); err != nil {
		if errors.Is(err, io.EOF) {
			return chatapp.ErrSessionClosed
		}
		if isDisconnectError(err) {
			return chatapp.ErrSessionClosed
		}

		return err
	}

	return nil
}

func (s *Session) Recv(ctx context.Context) (chat.Message, error) {
	if err := ctx.Err(); err != nil {
		return chat.Message{}, err
	}

	envelope, err := s.stream.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return chat.Message{}, chatapp.ErrSessionClosed
		}
		if isDisconnectError(err) {
			return chat.Message{}, chatapp.ErrSessionClosed
		}

		return chat.Message{}, err
	}

	return protocolchatv1.FromEnvelope(envelope), nil
}

func isDisconnectError(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() {
	case codes.Unavailable, codes.Canceled:
		return true
	default:
		return false
	}
}

func (s *Session) Close() error {
	return s.closeFn()
}
