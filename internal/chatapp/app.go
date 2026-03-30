package chatapp

import "context"

type StatusReporter interface {
	PrintStatus(format string, args ...any)
	PrintError(format string, args ...any)
}

type App struct {
	cfg       Config
	ui        StatusReporter
	transport Transport
}

func New(cfg Config, ui StatusReporter, transport Transport) *App {
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

func (a *App) runClient(ctx context.Context) (retErr error) {
	session, err := a.transport.Dial(ctx, a.cfg.PeerAddr)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := session.Close()
		if closeErr != nil && retErr == nil {
			retErr = closeErr
		}
	}()

	a.ui.PrintStatus("chat client connected to %s as %s", a.cfg.PeerAddr, a.cfg.Name)

	return nil
}

func (a *App) runServer(ctx context.Context) error {
	_, err := a.transport.Listen(ctx, a.cfg.ListenAddr)
	if err != nil {
		return err
	}

	a.ui.PrintStatus("chat server listening on %s as %s", a.cfg.ListenAddr, a.cfg.Name)

	return nil
}
