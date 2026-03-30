package chatv1

import (
	"testing"
	"time"

	apichatv1 "se-xp-2026-chat/api/chat/v1"
	"se-xp-2026-chat/internal/chat"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestToEnvelopePreservesMessageFields(t *testing.T) {
	t.Parallel()

	sentAt := time.Date(2026, time.March, 30, 12, 34, 56, 0, time.UTC)
	msg := chat.Message{
		Sender: "Alice",
		SentAt: sentAt,
		Text:   "hello",
	}

	envelope := ToEnvelope(msg)

	if envelope.GetSender() != "Alice" {
		t.Fatalf("unexpected sender: %q", envelope.GetSender())
	}
	if envelope.GetText() != "hello" {
		t.Fatalf("unexpected text: %q", envelope.GetText())
	}
	if got := envelope.GetSentAt().AsTime(); !got.Equal(sentAt) {
		t.Fatalf("unexpected sent_at: %v", got)
	}
}

func TestFromEnvelopePreservesMessageFields(t *testing.T) {
	t.Parallel()

	sentAt := time.Date(2026, time.March, 30, 12, 34, 56, 0, time.UTC)
	envelope := &apichatv1.Envelope{
		Sender: "Bob",
		SentAt: timestamppb.New(sentAt),
		Text:   "hi",
	}

	msg := FromEnvelope(envelope)

	if msg.Sender != "Bob" {
		t.Fatalf("unexpected sender: %q", msg.Sender)
	}
	if msg.Text != "hi" {
		t.Fatalf("unexpected text: %q", msg.Text)
	}
	if !msg.SentAt.Equal(sentAt) {
		t.Fatalf("unexpected sent_at: %v", msg.SentAt)
	}
}

func TestFromEnvelopeNilEnvelope(t *testing.T) {
	t.Parallel()

	msg := FromEnvelope(nil)

	if msg != (chat.Message{}) {
		t.Fatalf("expected zero message, got %#v", msg)
	}
}
