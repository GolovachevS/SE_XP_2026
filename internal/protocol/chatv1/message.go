package chatv1

import (
	"se-xp-2026-chat/api/chat/v1"
	"se-xp-2026-chat/internal/chat"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToEnvelope converts the internal message model into the protobuf transport shape.
func ToEnvelope(msg chat.Message) *chatv1.Envelope {
	return &chatv1.Envelope{
		Sender: msg.Sender,
		SentAt: timestamppb.New(msg.SentAt),
		Text:   msg.Text,
	}
}

// FromEnvelope converts a protobuf message received over gRPC into the domain model.
func FromEnvelope(envelope *chatv1.Envelope) chat.Message {
	if envelope == nil {
		return chat.Message{}
	}

	var sentAtValue = envelope.GetSentAt()

	msg := chat.Message{
		Sender: envelope.GetSender(),
		Text:   envelope.GetText(),
	}

	if sentAtValue != nil {
		msg.SentAt = sentAtValue.AsTime()
	}

	return msg
}
