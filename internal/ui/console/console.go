package console

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"se-xp-2026-chat/internal/chat"
)

// UI implements basic console input/output for the chat application.
type UI struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

// New creates UI bound to OS stdin/stdout/stderr.
func New(out, err io.Writer) *UI {
	return &UI{
		in:  os.Stdin,
		out: out,
		err: err,
	}
}

// NewWithStreams creates UI with explicit input/output streams.
// It is primarily used by tests.
func NewWithStreams(in io.Reader, out io.Writer, err io.Writer) *UI {
	return &UI{
		in:  in,
		out: out,
		err: err,
	}
}

func (ui *UI) PrintStatus(format string, args ...any) {
	_, _ = fmt.Fprintf(ui.out, format+"\n", args...)
}

func (ui *UI) PrintError(format string, args ...any) {
	_, _ = fmt.Fprintf(ui.err, format+"\n", args...)
}

// PrintMessage outputs an incoming or outgoing chat message in a unified format.
func (ui *UI) PrintMessage(msg chat.Message) {
	_, _ = fmt.Fprintln(ui.out, msg.Format())
}

// ReadLines reads user input line-by-line from stdin.
//
// It trims whitespace and filters out empty lines.
// The returned channel is closed on EOF.
//
// Note: cancellation can stop the consumer side, but may not be able to interrupt
// a blocking read from OS stdin.
func (ui *UI) ReadLines(ctx context.Context) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		scanner := bufio.NewScanner(ui.in)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line := chat.NormalizeText(scanner.Text())
			if chat.IsEmptyText(line) {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case out <- line:
			}
		}

		if err := scanner.Err(); err != nil {
			ui.PrintError("stdin read error: %v", err)
		}
	}()

	return out
}
