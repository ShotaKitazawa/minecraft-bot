package mock

import "errors"

type BotSenderMockValid struct {
	Msg string
}

func (sender *BotSenderMockValid) SendTextMessage(msg string) error {
	sender.Msg = msg
	return nil
}

type BotSenderMockInvalid struct{}

func (sender *BotSenderMockInvalid) SendTextMessage(msg string) error {
	return errors.New(``)
}
