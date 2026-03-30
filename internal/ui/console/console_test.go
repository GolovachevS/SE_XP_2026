package console

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"se-xp-2026-chat/internal/chat"
)

func TestPrintStatusWritesToStdout(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	var err bytes.Buffer
	ui := New(&out, &err)

	ui.PrintStatus("hello %s", "chat")

	if got := out.String(); got != "hello chat\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}

	if got := err.String(); got != "" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}

func TestPrintErrorWritesToStderr(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	var err bytes.Buffer
	ui := New(&out, &err)

	ui.PrintError("boom: %d", 42)

	if got := err.String(); got != "boom: 42\n" {
		t.Fatalf("unexpected stderr: %q", got)
	}

	if got := out.String(); got != "" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}

func TestReadLines_TrimsAndFiltersEmpty(t *testing.T) {
	t.Parallel()

	in := strings.NewReader("\n  \n  hello  \n\t\nworld\n")
	var out bytes.Buffer
	var err bytes.Buffer
	ui := NewWithStreams(in, &out, &err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var got []string
	for line := range ui.ReadLines(ctx) {
		got = append(got, line)
	}

	if len(got) != 2 || got[0] != "hello" || got[1] != "world" {
		t.Fatalf("unexpected lines: %#v", got)
	}
}

func TestPrintMessage_UsesDomainFormat(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	var err bytes.Buffer
	ui := NewWithStreams(strings.NewReader(""), &out, &err)

	sentAt := time.Date(2026, time.March, 30, 12, 0, 0, 0, time.UTC)
	msg, msgErr := chat.NewMessage("Alice", sentAt, "hello")
	if msgErr != nil {
		t.Fatalf("NewMessage returned error: %v", msgErr)
	}

	ui.PrintMessage(msg)

	if got := out.String(); got != msg.Format()+"\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}

	if got := err.String(); got != "" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}
