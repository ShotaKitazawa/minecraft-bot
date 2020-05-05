package line

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/sirupsen/logrus"
)

type BotAdaptor struct {
	botplug.BaseConfig
	ChannelSecret    string
	ChannelToken     string
	NotificationMode string
	Endpoint         string
	GroupIDs         []string
	bot              *linebot.Client
}

func New(logger *logrus.Logger, endpoint, channelSecret, channelToken, groupIDsStr, notificationMode string) (*BotAdaptor, error) {
	botplugLINEConfig := botplug.New(logger)
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}

	ba := &BotAdaptor{
		BaseConfig:       botplugLINEConfig,
		ChannelSecret:    channelSecret,
		ChannelToken:     channelToken,
		NotificationMode: notificationMode,
		Endpoint:         endpoint,
		GroupIDs:         strings.Split(groupIDsStr, ","),
		bot:              bot,
	}
	go ba.pushMessageLoop()
	return ba, nil
}

func connect(channelSecret, channelToken string) (*linebot.Client, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

func (bot *BotAdaptor) WithPlugin(plugin botplug.BotPlugin) *BotAdaptor {
	result := bot
	result.Plugins = append(result.Plugins, plugin)
	return bot
}

// SendTextMessage is pkg/botplug.BotSender interface's implementation
func (ba *BotAdaptor) SendTextMessageToChannels(text string) (err error) {
	for _, groupID := range ba.GroupIDs {
		if err := ba.sendTextMessage(groupID, text); err != nil {
			return err
		}
	}
	return nil
}
func (ba *BotAdaptor) sendTextMessage(groupID, text string) (err error) {
	_, err = ba.bot.PushMessage(groupID, linebot.NewTextMessage(text)).Do()
	return err
}
func (ba *BotAdaptor) replyTextMessage(replyToken, text string) (err error) {
	_, err = ba.bot.ReplyMessage(replyToken, linebot.NewTextMessage(text)).Do()
	return err
}

// pushMessageLoop is called when Bot instance created
func (ba *BotAdaptor) pushMessageLoop() (err error) {
	for {
		for _, plugin := range ba.Plugins {
			// execute user function
			output := plugin.PushMessageEntry()
			// proceed contents in queue
			texts, err := botplug.FormatToText(output)
			if err != nil {
				return err
			}
			for _, text := range texts {
				if text == "" {
					continue
				}
				if err := ba.SendTextMessageToChannels(text); err != nil {
					ba.Logger.Error(`failed to push notification: `, err)
				}
			}
		}
		time.Sleep(time.Millisecond)
	}
}

// call Bot Plugin interface
func (ba *BotAdaptor) NewHandler() (*httphandler.WebhookHandler, error) {
	handler, err := httphandler.New(
		ba.ChannelSecret,
		ba.ChannelToken,
	)
	if err != nil {
		return nil, err
	}
	ba.bot, err = connect(ba.ChannelSecret, ba.ChannelToken)
	if err != nil {
		return nil, err
	}

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				switch event.Message.(type) {
				case *linebot.TextMessage:
					if err = ba.receiveTextMessage(event); err != nil {
						log.Print(err)
					}
				}
			case linebot.EventTypeMemberJoined:
				if err = ba.receiveMemberJoin(event); err != nil {
					log.Print(err)
				}
			}
		}
	})
	return handler, nil
}

// receiveTextMessage is called when Receive Chat Message
func (receiver *BotAdaptor) receiveTextMessage(event *linebot.Event) (err error) {
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

	for _, plugin := range receiver.Plugins {
		// execute user function
		output := plugin.ReceiveMessageEntry(input)
		if output == nil {
			return
		}
		// proceed contents in queue
		texts, err := botplug.FormatToText(output)
		if err != nil {
			return err
		}
		for _, text := range texts {
			if err := receiver.replyTextMessage(event.ReplyToken, text); err != nil {
				return err
			}
		}
	}

	return nil
}

// receiveMemberJoin is called when someone join at Chat Group
func (receiver *BotAdaptor) receiveMemberJoin(event *linebot.Event) (err error) {
	input := &botplug.MessageInput{
		Timestamp: event.Timestamp,
		Source: &botplug.Source{
			Type:    string(event.Source.Type),
			UserID:  event.Source.UserID,
			GroupID: event.Source.GroupID,
		},
	}
	for _, plugin := range receiver.Plugins {
		// execute user function
		output := plugin.ReceiveMemberJoinEntry(input)
		if output == nil {
			return
		}
		// proceed contents in queue
		texts, err := botplug.FormatToText(output)
		if err != nil {
			return err
		}
		for _, text := range texts {
			if err := receiver.replyTextMessage(event.ReplyToken, text); err != nil {
				return err
			}
		}
	}
	return nil
}
