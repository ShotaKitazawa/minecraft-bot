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
func (sender *BotAdaptor) SendTextMessage(text string) (err error) {
	for _, groupID := range sender.GroupIDs {
		if _, err := sender.bot.PushMessage(groupID, linebot.NewTextMessage(text)).Do(); err != nil {
			sender.Logger.Error(`failed to push notification: `, err)
			return err
		}
	}
	return nil
}

// pushMessageLoop is called when Bot instance created
func (ba *BotAdaptor) pushMessageLoop() (err error) {
	for {
		for _, plugin := range ba.Plugins {
			// execute user function
			output := plugin.PushMessageEntry()
			// proceed contents in queue
			if err := ba.pushFromQueue(output); err != nil {
				ba.Logger.Error(err)
			}
		}
		time.Sleep(time.Second)
	}
}

func (ba *BotAdaptor) pushFromQueue(output *botplug.MessageOutput) (err error) {
	// proceed contents in queue
	for _, element := range output.Queue {
		switch typedElement := element.(type) {
		case string:
			if typedElement == "" {
				return
			}
			if err = ba.SendTextMessage(typedElement); err != nil {
				return
			}
		case []string:
			if len(typedElement) == 0 {
				return
			}
			if err = ba.SendTextMessage(strings.Join(typedElement, ",")); err != nil {
				return
			}
		case error:
			if typedElement.Error() == "" {
				return
			}
			if err = ba.SendTextMessage(typedElement.Error()); err != nil {
				return
			}
		case []linebot.SendingMessage:
			// TODO: implement
		}
	}
	return nil
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
		if err := receiver.replyFromQueue(event, output); err != nil {
			return err
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
		if err := receiver.replyFromQueue(event, output); err != nil {
			return err
		}

	}
	return nil
}

func (receiver *BotAdaptor) replyFromQueue(event *linebot.Event, output *botplug.MessageOutput) (err error) {
	// proceed contents in queue
	for _, element := range output.Queue {
		switch typedElement := element.(type) {
		case string:
			if _, err = receiver.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(typedElement)).Do(); err != nil {
				return
			}
		case []string:
			if _, err = receiver.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(strings.Join(typedElement, ","))).Do(); err != nil {
				return
			}
		case error:
			if _, err = receiver.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(typedElement.Error())).Do(); err != nil {
				return
			}
		case []linebot.SendingMessage:
			if _, err = receiver.bot.ReplyMessage(event.ReplyToken, typedElement...).Do(); err != nil {
				return
			}
		}
	}
	return nil
}
