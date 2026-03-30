package grpcchat

import (
	"context"
	"errors"
	"net"

	apichatv1 "se-xp-2026-chat/api/chat/v1"
	chatapp "se-xp-2026-chat/internal/app"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Transport struct{}

func New() *Transport {
	return &Transport{}
}

func (t *Transport) Dial(ctx context.Context, peerAddr string) (chatapp.Session, error) {
	conn, err := grpc.NewClient(
		peerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	stream, err := apichatv1.NewChatServiceClient(conn).ChatStream(ctx)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &Session{
		stream: stream,
		closeFn: func() error {
			closeErr := stream.CloseSend()
			connErr := conn.Close()

			if closeErr != nil {
				return closeErr
			}

			return connErr
		},
	}, nil
}

func (t *Transport) Listen(ctx context.Context, addr string) (<-chan chatapp.Session, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	sessions := make(chan chatapp.Session, 1)
	server := grpc.NewServer()
	apichatv1.RegisterChatServiceServer(server, chatServer{sessions: sessions})

	go func() {
		<-ctx.Done()
		server.GracefulStop()
		_ = listener.Close()
	}()

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, net.ErrClosed) {
			server.Stop()
		}
	}()

	return sessions, nil
}

type chatServer struct {
	apichatv1.UnimplementedChatServiceServer
	sessions chan<- chatapp.Session
}

func (s chatServer) ChatStream(stream apichatv1.ChatService_ChatStreamServer) error {
	session := &Session{
		stream:  stream,
		closeFn: func() error { return nil },
	}

	select {
	case s.sessions <- session:
	case <-stream.Context().Done():
	}

	<-stream.Context().Done()

	if errors.Is(stream.Context().Err(), context.Canceled) {
		return nil
	}

	return stream.Context().Err()
}
