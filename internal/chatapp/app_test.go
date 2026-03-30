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

	sessions := make(chan Session, 1)
	sessions <- &stubSession{}
	close(sessions)

	transport := &stubTransport{sessions: sessions}
	ui := &stubUI{}
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

	if got := ui.statuses[0]; got != "chat server listening on :50051 as Alice" {
		t.Fatalf("unexpected status: %q", got)
	}
}

func TestAppRunClientUsesDial(t *testing.T) {
	t.Parallel()

	transport := &stubTransport{}
	ui := &stubUI{}
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

	if len(ui.statuses) != 1 {
		t.Fatalf("unexpected statuses count: %d", len(ui.statuses))
	}

	if got := ui.statuses[0]; got != "chat client connected to 127.0.0.1:50051 as Bob" {
		t.Fatalf("unexpected status: %q", got)
	}
}

func TestAppRunPropagatesListenError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("listen failed")
	app := New(
		Config{Name: "Alice", ListenAddr: ":50051"},
		&stubUI{},
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
		&stubUI{},
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
		&stubUI{},
		transport,
	)

	err := app.Run(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

type stubUI struct {
	statuses []string
	errors   []string
	lines    []string
	msgs     []chat.Message
}

func (s *stubUI) PrintStatus(format string, args ...any) {
	s.statuses = append(s.statuses, fmt.Sprintf(format, args...))
}

func (s *stubUI) PrintError(format string, args ...any) {
	s.errors = append(s.errors, fmt.Sprintf(format, args...))
}

func (s *stubUI) PrintMessage(msg chat.Message) {
	s.msgs = append(s.msgs, msg)
}

func (s *stubUI) ReadLines(context.Context) <-chan string {
	ch := make(chan string, len(s.lines))
	for _, line := range s.lines {
		ch <- line
	}
	close(ch)
	return ch
}

type stubTransport struct {
	listenAddr string
	dialAddr   string
	listenErr  error
	dialErr    error
	session    stubSession
	sessions   <-chan Session
}

func (s *stubTransport) Listen(_ context.Context, addr string) (<-chan Session, error) {
	if s.listenErr != nil {
		return nil, s.listenErr
	}

	s.listenAddr = addr
	if s.sessions != nil {
		return s.sessions, nil
	}

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
