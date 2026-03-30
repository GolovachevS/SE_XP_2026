package chat

import (
	"testing"
	"time"
)

func TestNewMessageTrimsFields(t *testing.T) {
	t.Parallel()

	sentAt := time.Date(2026, time.March, 30, 12, 0, 0, 0, time.UTC)

	msg, err := NewMessage(" Alice ", sentAt, " hello ")
	if err != nil {
		t.Fatalf("NewMessage returned error: %v", err)
	}

	if msg.Sender != "Alice" {
		t.Fatalf("unexpected sender: %q", msg.Sender)
	}

	if msg.Text != "hello" {
		t.Fatalf("unexpected text: %q", msg.Text)
	}

	if !msg.SentAt.Equal(sentAt) {
		t.Fatalf("unexpected sent_at: %v", msg.SentAt)
	}
}

func TestMessageValidateRejectsMissingSender(t *testing.T) {
	t.Parallel()

	msg := Message{
		Sender: "   ",
		SentAt: time.Now(),
		Text:   "hello",
	}

	if err := msg.Validate(); err != ErrSenderRequired {
		t.Fatalf("expected ErrSenderRequired, got %v", err)
	}
}

func TestMessageValidateRejectsMissingTimestamp(t *testing.T) {
	t.Parallel()

	msg := Message{
		Sender: "Alice",
		Text:   "hello",
	}

	if err := msg.Validate(); err != ErrSentAtRequired {
		t.Fatalf("expected ErrSentAtRequired, got %v", err)
	}
}

func TestMessageValidateRejectsMissingText(t *testing.T) {
	t.Parallel()

	msg := Message{
		Sender: "Alice",
		SentAt: time.Now(),
		Text:   "   ",
	}

	if err := msg.Validate(); err != ErrTextRequired {
		t.Fatalf("expected ErrTextRequired, got %v", err)
	}
}

func TestMessageFormat_IsStable(t *testing.T) {
	t.Parallel()

	sentAt := time.Date(2026, time.March, 30, 12, 0, 0, 0, time.UTC)
	msg, err := NewMessage("Alice", sentAt, "hello")
	if err != nil {
		t.Fatalf("NewMessage returned error: %v", err)
	}

	want := "2026-03-30 12:00:00 Alice: hello"
	if got := msg.Format(); got != want {
		t.Fatalf("unexpected format: %q", got)
	}
}

func TestNormalizeText_TrimsWhitespace(t *testing.T) {
	t.Parallel()

	if got := NormalizeText("  hi\t"); got != "hi" {
		t.Fatalf("unexpected normalized text: %q", got)
	}
}

func TestIsEmptyText_AfterNormalization(t *testing.T) {
	t.Parallel()

	if !IsEmptyText("  \n\t  ") {
		t.Fatalf("expected empty text")
	}
}
