package chatapp

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"se-xp-2026-chat/internal/chat"
)

func TestAppRunServerUsesListen(t *testing.T) {
	t.Parallel()

	transport := &stubTransport{}
	ui := &stubReporter{}
	app := New(Config{
		Name:       "Alice",
		ListenAddr: ":50051",
	}, ui, transport)

	if err := app.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if transport.listenAddr != ":50051" {
		t.Fatalf("unexpected listen addr: %q", transport.listenAddr)
	}

	if transport.dialAddr != "" {
		t.Fatalf("unexpected dial addr: %q", transport.dialAddr)
	}

	if len(ui.statuses) != 1 {
		t.Fatalf("unexpected statuses count: %d", len(ui.statuses))
	}
}

func TestAppRunClientUsesDial(t *testing.T) {
	t.Parallel()

	transport := &stubTransport{}
	ui := &stubReporter{}
	app := New(Config{
		Name:     "Bob",
		PeerAddr: "127.0.0.1:50051",
	}, ui, transport)

	if err := app.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if transport.dialAddr != "127.0.0.1:50051" {
		t.Fatalf("unexpected dial addr: %q", transport.dialAddr)
	}

	if transport.listenAddr != "" {
		t.Fatalf("unexpected listen addr: %q", transport.listenAddr)
	}

	if transport.session.closeCalls != 1 {
		t.Fatalf("expected session to be closed once, got %d", transport.session.closeCalls)
	}
}

func TestAppRunPropagatesListenError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("listen failed")
	app := New(
		Config{Name: "Alice", ListenAddr: ":50051"},
		&stubReporter{},
		&stubTransport{listenErr: wantErr},
	)

	err := app.Run(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestAppRunPropagatesDialError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("dial failed")
	app := New(
		Config{Name: "Bob", PeerAddr: "127.0.0.1:50051"},
		&stubReporter{},
		&stubTransport{dialErr: wantErr},
	)

	err := app.Run(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestAppRunPropagatesCloseError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("close failed")
	transport := &stubTransport{
		session: stubSession{closeErr: wantErr},
	}
	app := New(
		Config{Name: "Bob", PeerAddr: "127.0.0.1:50051"},
		&stubReporter{},
		transport,
	)

	err := app.Run(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

type stubReporter struct {
	statuses []string
	errors   []string
}

func (s *stubReporter) PrintStatus(format string, args ...any) {
	s.statuses = append(s.statuses, fmt.Sprintf(format, args...))
}

func (s *stubReporter) PrintError(format string, args ...any) {
	s.errors = append(s.errors, fmt.Sprintf(format, args...))
}

type stubTransport struct {
	listenAddr string
	dialAddr   string
	listenErr  error
	dialErr    error
	session    stubSession
}

func (s *stubTransport) Listen(_ context.Context, addr string) (<-chan Session, error) {
	if s.listenErr != nil {
		return nil, s.listenErr
	}

	s.listenAddr = addr

	return make(chan Session), nil
}

func (s *stubTransport) Dial(_ context.Context, peerAddr string) (Session, error) {
	if s.dialErr != nil {
		return nil, s.dialErr
	}

	s.dialAddr = peerAddr

	return &s.session, nil
}

type stubSession struct {
	closeCalls int
	closeErr   error
}

func (*stubSession) Send(context.Context, chat.Message) error {
	return nil
}

func (*stubSession) Recv(context.Context) (chat.Message, error) {
	return chat.Message{}, ErrSessionClosed
}

func (s *stubSession) Close() error {
	s.closeCalls++
	return s.closeErr
}
