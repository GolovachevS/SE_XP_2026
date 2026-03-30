package console

import (
	"bytes"
	"context"
	"errors"
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

func TestReadLines_ReportsScannerError(t *testing.T) {
	t.Parallel()

	in := &errorAfterReader{
		data: []byte("hello\n"),
		err:  errors.New("boom"),
	}
	var out bytes.Buffer
	var err bytes.Buffer
	ui := NewWithStreams(in, &out, &err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var got []string
	for line := range ui.ReadLines(ctx) {
		got = append(got, line)
	}

	if len(got) != 1 || got[0] != "hello" {
		t.Fatalf("unexpected lines: %#v", got)
	}

	if stderr := err.String(); !strings.Contains(stderr, "stdin read error") || !strings.Contains(stderr, "boom") {
		t.Fatalf("expected read error to be reported, got: %q", stderr)
	}
}

type errorAfterReader struct {
	data []byte
	pos  int
	err  error
}

func (r *errorAfterReader) Read(p []byte) (int, error) {
	if r.pos < len(r.data) {
		n := copy(p, r.data[r.pos:])
		r.pos += n
		return n, nil
	}

	return 0, r.err
}
