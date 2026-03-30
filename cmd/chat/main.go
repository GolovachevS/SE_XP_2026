package main

import (
	"context"
	"os"

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

	if err := app.Run(context.Background()); err != nil {
		ui.PrintError("application error: %v", err)
		os.Exit(1)
	}
}
