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

type Message struct {
	Sender string
	SentAt time.Time
	Text   string
}

func NewMessage(sender string, sentAt time.Time, text string) (Message, error) {
	msg := Message{
		Sender: strings.TrimSpace(sender),
		SentAt: sentAt,
		Text:   strings.TrimSpace(text),
	}

	if err := msg.Validate(); err != nil {
		return Message{}, err
	}

	return msg, nil
}

func (m Message) Validate() error {
	if strings.TrimSpace(m.Sender) == "" {
		return ErrSenderRequired
	}

	if m.SentAt.IsZero() {
		return ErrSentAtRequired
	}

	if strings.TrimSpace(m.Text) == "" {
		return ErrTextRequired
	}

	return nil
}
