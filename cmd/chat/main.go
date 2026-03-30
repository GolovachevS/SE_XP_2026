package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"se-xp-2026-chat/internal/chatapp"
	"se-xp-2026-chat/internal/transport/grpcchat"
	"se-xp-2026-chat/internal/ui/console"
)

func main() {
	cfg, err := chatapp.ParseArgs(os.Args[1:])
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	ui := console.New(os.Stdout, os.Stderr)
	app := chatapp.New(cfg, ui, grpcchat.New())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		ui.PrintError("application error: %v", err)
		os.Exit(1)
	}
}
