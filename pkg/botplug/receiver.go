package botplug

import "time"

type BotPlugin interface {
	ReceiveMessageEntry(*MessageInput) *MessageOutput
	ReceiveMemberJoinEntry(*MessageInput) *MessageOutput
	PushMessageEntry() *MessageOutput
}

type MessageInput struct {
	Timestamp time.Time
	Source    *Source
	Messages  []string
}

type Source struct {
	Type    string
	UserID  string
	GroupID string
}

type MessageOutput struct {
	Queue []interface{}
}
