package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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

func TestAppRunServerReturnsContextErrorWhileWaitingForSession(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	transport := &stubTransport{sessions: make(chan Session)}
	ui := &stubUI{}
	app := New(
		Config{Name: "Alice", ListenAddr: ":50051"},
		ui,
		transport,
	)

	err := app.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestSendLoopSendsMessagesFromUI(t *testing.T) {
	t.Parallel()

	ui := &stubUI{lines: []string{"  hello  "}}
	session := &stubSession{}
	app := New(
		Config{Name: "Alice"},
		ui,
		&stubTransport{},
	)

	err := app.sendLoop(context.Background(), session)
	if err != nil {
		t.Fatalf("sendLoop returned error: %v", err)
	}

	if len(session.sent) != 1 {
		t.Fatalf("expected 1 sent message, got %d", len(session.sent))
	}

	got := session.sent[0]
	if got.Sender != "Alice" {
		t.Fatalf("unexpected sender: %q", got.Sender)
	}
	if got.Text != "hello" {
		t.Fatalf("unexpected text: %q", got.Text)
	}
	if got.SentAt.IsZero() {
		t.Fatal("expected non-zero sent time")
	}
}

func TestRecvLoopPrintsMessagesUntilSessionClosed(t *testing.T) {
	t.Parallel()

	msg := chat.Message{
		Sender: "Bob",
		SentAt: time.Date(2026, time.March, 30, 12, 0, 0, 0, time.UTC),
		Text:   "hi",
	}
	session := &stubSession{
		recvMsgs: []chat.Message{msg},
	}
	ui := &stubUI{}
	app := New(
		Config{Name: "Alice"},
		ui,
		&stubTransport{},
	)

	err := app.recvLoop(context.Background(), session)
	if err != nil {
		t.Fatalf("recvLoop returned error: %v", err)
	}

	if len(ui.msgs) != 1 {
		t.Fatalf("expected 1 printed message, got %d", len(ui.msgs))
	}

	if got := ui.msgs[0]; got != msg {
		t.Fatalf("unexpected printed message: %#v", got)
	}
}

func TestSendLoopPropagatesSessionError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("send failed")
	ui := &stubUI{lines: []string{"hello"}}
	session := &stubSession{sendErr: wantErr}
	app := New(
		Config{Name: "Alice"},
		ui,
		&stubTransport{},
	)

	err := app.sendLoop(context.Background(), session)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestRecvLoopPropagatesUnexpectedError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("recv failed")
	session := &stubSession{recvErr: wantErr}
	app := New(
		Config{Name: "Alice"},
		&stubUI{},
		&stubTransport{},
	)

	err := app.recvLoop(context.Background(), session)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestRunSessionReturnsNilOnRemoteDisconnectWhileInputWaits(t *testing.T) {
	t.Parallel()

	ui := &blockingUI{}
	session := &stubSession{recvErr: ErrSessionClosed}
	app := New(
		Config{Name: "Alice"},
		ui,
		&stubTransport{},
	)

	err := app.runSession(context.Background(), session)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if session.closeCalls != 1 {
		t.Fatalf("expected session to be closed once, got %d", session.closeCalls)
	}
}

func TestRunSessionStopsOnEOFFromInput(t *testing.T) {
	t.Parallel()

	ui := &stubUI{}
	session := &stubSession{}
	app := New(
		Config{Name: "Alice"},
		ui,
		&stubTransport{},
	)

	err := app.runSession(context.Background(), session)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if session.closeCalls != 1 {
		t.Fatalf("expected session to be closed once, got %d", session.closeCalls)
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

type blockingUI struct {
	stubUI
}

func (b *blockingUI) ReadLines(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		<-ctx.Done()
		close(ch)
	}()
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
	sendErr    error
	sent       []chat.Message
	recvMsgs   []chat.Message
	recvErr    error
	recvIndex  int
}

func (s *stubSession) Send(_ context.Context, msg chat.Message) error {
	if s.sendErr != nil {
		return s.sendErr
	}

	s.sent = append(s.sent, msg)
	return nil
}

func (s *stubSession) Recv(context.Context) (chat.Message, error) {
	if s.recvIndex < len(s.recvMsgs) {
		msg := s.recvMsgs[s.recvIndex]
		s.recvIndex++
		return msg, nil
	}

	if s.recvErr != nil {
		return chat.Message{}, s.recvErr
	}

	return chat.Message{}, ErrSessionClosed
}

func (s *stubSession) Close() error {
	s.closeCalls++
	return s.closeErr
}
