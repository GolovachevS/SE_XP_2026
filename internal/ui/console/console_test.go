package console

import (
	"bytes"
	"testing"
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
