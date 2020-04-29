package line

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

type BotSender struct {
	Config
	GroupIDs []string
}

func (sender *BotSender) SendTextMessage(text string) (err error) {
	for _, groupID := range sender.GroupIDs {
		if _, err := sender.bot.PushMessage(groupID, linebot.NewTextMessage(text)).Do(); err != nil {
			sender.Logger.Error(`failed to push notification: `, err)
			return err
		}
	}
	return nil
}
