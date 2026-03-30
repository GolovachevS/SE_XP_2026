package grpcchat

import (
	"context"
	"net"
	"testing"
	"time"

	apichatv1 "se-xp-2026-chat/api/chat/v1"
	chatapp "se-xp-2026-chat/internal/app"
	"se-xp-2026-chat/internal/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestChatStreamAcceptsServerSessionOverBufconn(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listener, sessions, cleanup := startBufconnServer(t)
	defer cleanup()

	clientSession, clientCleanup := newBufconnClientSession(t, ctx, listener)
	defer clientCleanup()
	defer func() { _ = clientSession.Close() }()

	serverSession := waitForAcceptedSession(t, ctx, sessions)
	if serverSession == nil {
		t.Fatal("expected accepted server session")
	}
}

func TestChatStreamExchangesMessagesOverBufconn(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listener, sessions, cleanup := startBufconnServer(t)
	defer cleanup()

	clientSession, clientCleanup := newBufconnClientSession(t, ctx, listener)
	defer clientCleanup()
	defer func() { _ = clientSession.Close() }()

	serverSession := waitForAcceptedSession(t, ctx, sessions)

	clientMsg := chat.Message{
		Sender: "Alice",
		SentAt: time.Date(2026, time.March, 30, 12, 0, 0, 0, time.UTC),
		Text:   "hello",
	}
	if err := clientSession.Send(ctx, clientMsg); err != nil {
		t.Fatalf("client send returned error: %v", err)
	}

	gotOnServer, err := serverSession.Recv(ctx)
	if err != nil {
		t.Fatalf("server recv returned error: %v", err)
	}
	if gotOnServer != clientMsg {
		t.Fatalf("unexpected server message: %#v", gotOnServer)
	}

	serverMsg := chat.Message{
		Sender: "Bob",
		SentAt: time.Date(2026, time.March, 30, 12, 1, 0, 0, time.UTC),
		Text:   "hi",
	}
	if err := serverSession.Send(ctx, serverMsg); err != nil {
		t.Fatalf("server send returned error: %v", err)
	}

	gotOnClient, err := clientSession.Recv(ctx)
	if err != nil {
		t.Fatalf("client recv returned error: %v", err)
	}
	if gotOnClient != serverMsg {
		t.Fatalf("unexpected client message: %#v", gotOnClient)
	}
}

func startBufconnServer(t *testing.T) (*bufconn.Listener, <-chan chatapp.Session, func()) {
	t.Helper()

	listener := bufconn.Listen(1024 * 1024)
	sessions := make(chan chatapp.Session, 1)
	server := grpc.NewServer()
	apichatv1.RegisterChatServiceServer(server, chatServer{sessions: sessions})

	go func() {
		_ = server.Serve(listener)
	}()

	cleanup := func() {
		server.Stop()
		_ = listener.Close()
	}

	return listener, sessions, cleanup
}

func newBufconnClientSession(t *testing.T, ctx context.Context, listener *bufconn.Listener) (*Session, func()) {
	t.Helper()

	conn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient returned error: %v", err)
	}

	stream, err := apichatv1.NewChatServiceClient(conn).ChatStream(ctx)
	if err != nil {
		_ = conn.Close()
		t.Fatalf("ChatStream returned error: %v", err)
	}

	cleanup := func() {
		_ = stream.CloseSend()
		_ = conn.Close()
	}

	return NewClientSession(stream), cleanup
}

func waitForAcceptedSession(t *testing.T, ctx context.Context, sessions <-chan chatapp.Session) chatapp.Session {
	t.Helper()

	select {
	case <-ctx.Done():
		t.Fatalf("timed out waiting for accepted session: %v", ctx.Err())
	case session := <-sessions:
		return session
	}

	return nil
}
