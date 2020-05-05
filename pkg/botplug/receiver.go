package botplug

import (
	"strings"
	"time"
)

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

func FormatToText(output *MessageOutput) (result []string, err error) {
	for _, element := range output.Queue {
		switch typedElement := element.(type) {
		case string:
			result = append(result, typedElement)
		case []string:
			result = append(result, strings.Join(typedElement, ","))
		case error:
			result = append(result, typedElement.Error())
		}
	}
	return result, nil
}
