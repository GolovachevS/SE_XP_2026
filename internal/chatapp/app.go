package chatapp

import "context"

type StatusReporter interface {
	PrintStatus(format string, args ...any)
	PrintError(format string, args ...any)
}

type App struct {
	cfg Config
	ui  StatusReporter
}

func New(cfg Config, ui StatusReporter) *App {
	return &App{
		cfg: cfg,
		ui:  ui,
	}
}

func (a *App) Run(_ context.Context) error {
	a.ui.PrintStatus(
		"chat skeleton started: mode=%s name=%s listen=%s peer=%s",
		a.cfg.Mode(),
		a.cfg.Name,
		a.cfg.ListenAddr,
		a.cfg.PeerAddr,
	)

	return nil
}
