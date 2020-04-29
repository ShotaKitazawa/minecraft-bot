package line

import (
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
)

type BotReceiver struct {
	Config
	Plugin botplug.BotPluginReceiver
}

func (receiver *BotReceiver) WithPlugin(plugin botplug.BotPluginReceiver) *BotReceiver {
	result := receiver
	result.Plugin = plugin
	return result
}

func (receiver *BotReceiver) NewHandler() (*httphandler.WebhookHandler, error) {
	handler, err := httphandler.New(
		receiver.ChannelSecret,
		receiver.ChannelToken,
	)
	if err != nil {
		return nil, err
	}
	receiver.bot, err = connect(receiver.Config)
	if err != nil {
		return nil, err
	}

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				switch event.Message.(type) {
				case *linebot.TextMessage:
					if err = receiver.ReceiveTextMessage(event); err != nil {
						log.Print(err)
					}
				}
			case linebot.EventTypeMemberJoined:
				if err = receiver.ReceiveMemberJoin(event); err != nil {
					log.Print(err)
				}
			}
		}
	})
	return handler, nil
}

func (receiver *BotReceiver) ReceiveTextMessage(event *linebot.Event) (err error) {
	message := event.Message.(*linebot.TextMessage)
	input := &botplug.MessageInput{
		Timestamp: event.Timestamp,
		Source: &botplug.Source{
			Type:    string(event.Source.Type),
			UserID:  event.Source.UserID,
			GroupID: event.Source.GroupID,
		},
		Messages: strings.Fields(message.Text),
	}

	// execute user function
	output := receiver.Plugin.ReceiveMessageEntry(input)
	if output == nil {
		return
	}

	// proceed contents in queue
	if err := receiver.sendQueue(event, receiver.bot, output); err != nil {
		return err
	}

	return nil
}

func (receiver *BotReceiver) ReceiveMemberJoin(event *linebot.Event) (err error) {
	input := &botplug.MessageInput{
		Timestamp: event.Timestamp,
		Source: &botplug.Source{
			Type:    string(event.Source.Type),
			UserID:  event.Source.UserID,
			GroupID: event.Source.GroupID,
		},
	}

	// execute user function
	output := receiver.Plugin.ReceiveMemberJoinEntry(input)
	if output == nil {
		return
	}

	// proceed contents in queue
	if err := receiver.sendQueue(event, receiver.bot, output); err != nil {
		return err
	}

	return nil
}

func (receiver *BotReceiver) sendQueue(event *linebot.Event, bot *linebot.Client, output *botplug.MessageOutput) (err error) {
	// proceed contents in queue
	for _, element := range output.Queue {
		switch typedElement := element.(type) {
		case string:
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(typedElement)).Do(); err != nil {
				return
			}
		case []string:
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(strings.Join(typedElement, ","))).Do(); err != nil {
				return
			}
		case error:
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(typedElement.Error())).Do(); err != nil {
				return
			}
		case []linebot.SendingMessage:
			if _, err = bot.ReplyMessage(event.ReplyToken, typedElement...).Do(); err != nil {
				return
			}
		}
	}
	return nil
}
