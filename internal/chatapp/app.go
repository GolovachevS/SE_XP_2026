package chatapp

import (
	"context"
	"errors"
	"sync"
	"time"

	"se-xp-2026-chat/internal/chat"
)

// UI abstracts console input/output for orchestration.
//
// It is kept small on purpose to avoid coupling transport and rendering.
type UI interface {
	PrintStatus(format string, args ...any)
	PrintError(format string, args ...any)
	PrintMessage(msg chat.Message)
	ReadLines(ctx context.Context) <-chan string
}

type App struct {
	cfg       Config
	ui        UI
	transport Transport
}

func New(cfg Config, ui UI, transport Transport) *App {
	return &App{
		cfg:       cfg,
		ui:        ui,
		transport: transport,
	}
}

func (a *App) Run(ctx context.Context) error {
	if a.cfg.Mode() == "client" {
		return a.runClient(ctx)
	}

	return a.runServer(ctx)
}

func (a *App) runClient(ctx context.Context) error {
	session, err := a.transport.Dial(ctx, a.cfg.PeerAddr)
	if err != nil {
		return err
	}
	a.ui.PrintStatus("chat client connected to %s as %s", a.cfg.PeerAddr, a.cfg.Name)
	return a.runSession(ctx, session)
}

func (a *App) runServer(ctx context.Context) error {
	sessions, err := a.transport.Listen(ctx, a.cfg.ListenAddr)
	if err != nil {
		return err
	}

	a.ui.PrintStatus("chat server listening on %s as %s", a.cfg.ListenAddr, a.cfg.Name)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case session, ok := <-sessions:
		if !ok {
			return nil
		}
		return a.runSession(ctx, session)
	}
}

func (a *App) runSession(ctx context.Context, session Session) (retErr error) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var closeOnce sync.Once
	var closeErr error
	closeSession := func() {
		closeOnce.Do(func() {
			closeErr = session.Close()
		})
	}
	defer func() {
		closeSession()
		if retErr == nil && closeErr != nil {
			retErr = closeErr
		}
	}()

	recvDone := make(chan error, 1)
	go func() {
		recvDone <- a.recvLoop(childCtx, session)
	}()

	sendErr := a.sendLoop(childCtx, session)
	cancel()
	closeSession()

	recvErr := <-recvDone
	if sendErr != nil {
		return sendErr
	}
	if recvErr != nil {
		return recvErr
	}

	return nil
}

func (a *App) sendLoop(ctx context.Context, session Session) error {
	for line := range a.ui.ReadLines(ctx) {
		msg, err := chat.NewMessage(a.cfg.Name, time.Now(), line)
		if err != nil {
			a.ui.PrintError("invalid input: %v", err)
			continue
		}

		if err := session.Send(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) recvLoop(ctx context.Context, session Session) error {
	for {
		msg, err := session.Recv(ctx)
		if err != nil {
			if errors.Is(err, ErrSessionClosed) {
				return nil
			}
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		a.ui.PrintMessage(msg)
	}
}
