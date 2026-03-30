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

// Transport is the concrete gRPC-backed implementation of chatapp.Transport.
type Transport struct{}

// New creates a transport that dials and serves chat sessions over gRPC.
func New() *Transport {
	return &Transport{}
}

// Dial opens an outgoing bidirectional chat stream to the target peer.
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

// Listen starts a gRPC server and exposes accepted chat sessions through a channel.
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
		// GracefulStop lets active streams finish before the listener is closed.
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

// ChatStream hands the accepted server-side stream to the app layer as a Session.
func (s chatServer) ChatStream(stream apichatv1.ChatService_ChatStreamServer) error {
	session := &Session{
		stream:  stream,
		closeFn: func() error { return nil },
	}

	select {
	case s.sessions <- session:
	case <-stream.Context().Done():
	}

	// The app layer owns the session lifecycle after accepting it, so the RPC
	// handler only needs to stay alive until the stream context is canceled.
	<-stream.Context().Done()

	if errors.Is(stream.Context().Err(), context.Canceled) {
		return nil
	}

	return stream.Context().Err()
}
