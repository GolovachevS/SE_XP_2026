package chat

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrSenderRequired = errors.New("sender is required")
	ErrSentAtRequired = errors.New("sent_at is required")
	ErrTextRequired   = errors.New("text is required")
)

// TimeLayout is the canonical timestamp layout used in message rendering.
//
// Go's time formatting uses a reference time (Mon Jan 2 15:04:05 MST 2006).
// The chosen layout is human-readable and stable for tests.
const TimeLayout = "2006-01-02 15:04:05"

const senderSeparator = ": "

// NormalizeText applies basic normalization for user-provided text.
func NormalizeText(text string) string {
	return strings.TrimSpace(text)
}

// IsEmptyText reports whether text is empty after normalization.
func IsEmptyText(text string) bool {
	return NormalizeText(text) == ""
}

// Message is an immutable-ish chat domain object.
type Message struct {
	Sender string
	SentAt time.Time
	Text   string
}

// NewMessage constructs a validated Message.
func NewMessage(sender string, sentAt time.Time, text string) (Message, error) {
	msg := Message{
		Sender: NormalizeText(sender),
		SentAt: sentAt,
		Text:   NormalizeText(text),
	}

	if err := msg.Validate(); err != nil {
		return Message{}, err
	}

	return msg, nil
}

// Validate checks that a message contains all fields required for sending/rendering.
func (m Message) Validate() error {
	if IsEmptyText(m.Sender) {
		return ErrSenderRequired
	}

	if m.SentAt.IsZero() {
		return ErrSentAtRequired
	}

	if IsEmptyText(m.Text) {
		return ErrTextRequired
	}

	return nil
}

// Format renders the message for console output: time + sender + text.
func (m Message) Format() string {
	return m.SentAt.Format(TimeLayout) + " " + m.Sender + senderSeparator + m.Text
}
