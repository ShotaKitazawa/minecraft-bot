package slack

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
)

type BotAdaptor struct {
	botplug.BaseConfig
	Token            string
	NotificationMode string
	ChannelIDs       []string
	bot              *slack.Client
}

func New(logger *logrus.Logger, token, channelIDsStr, notificationMode string) (*BotAdaptor, error) {
	botplugSlackConfig := botplug.New(logger)
	bot := slack.New(token)

	ba := &BotAdaptor{
		BaseConfig:       botplugSlackConfig,
		Token:            token,
		NotificationMode: notificationMode,
		ChannelIDs:       strings.Split(channelIDsStr, ","),
		bot:              bot,
	}
	go ba.pushMessageLoop()
	return ba, nil
}

func (bot *BotAdaptor) WithPlugin(plugin botplug.BotPlugin) *BotAdaptor {
	result := bot
	result.Plugins = append(result.Plugins, plugin)
	return bot
}

// SendTextMessage is pkg/botplug.BotSender interface's implementation
func (sender *BotAdaptor) SendTextMessage(text string) (err error) {
	for _, channelID := range sender.ChannelIDs {
		if _, _, err := sender.bot.PostMessage(channelID, slack.MsgOptionText(text, false)); err != nil {
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
		}
	}
	return nil
}

// call Bot Plugin interface
func (ba *BotAdaptor) Run() error {
	rtm := ba.bot.NewRTM()
	go rtm.ManageConnection()
	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			for _, channelName := range ba.ChannelIDs {
				if ev.Channel == channelName {
					if err := ba.receiveTextMessage(ev); err != nil {
						ba.Logger.Error(err)
					}
				}
			}
		}
	}
	return nil
}

// receiveTextMessage is called when Receive Chat Message
func (receiver *BotAdaptor) receiveTextMessage(event *slack.MessageEvent) (err error) {
	input := &botplug.MessageInput{
		//TODO
		//Timestamp: event.Timestamp,
		Source: &botplug.Source{
			Type:    event.Type,
			UserID:  event.User,
			GroupID: event.Channel,
		},
		Messages: strings.Fields(event.Msg.Text),
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
func (receiver *BotAdaptor) receiveMemberJoin(event *slack.MessageEvent) (err error) {
	input := &botplug.MessageInput{
		//Timestamp: event.Timestamp,
		Source: &botplug.Source{
			Type:    string(event.Type),
			UserID:  event.User,
			GroupID: event.Channel,
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

func (receiver *BotAdaptor) replyFromQueue(event *slack.MessageEvent, output *botplug.MessageOutput) (err error) {
	// proceed contents in queue
	for _, element := range output.Queue {
		switch typedElement := element.(type) {
		case string:
			if _, _, err = receiver.bot.PostMessage(event.Channel, slack.MsgOptionText(typedElement, false)); err != nil {
				return
			}
		case []string:
			if _, _, err = receiver.bot.PostMessage(event.Channel, slack.MsgOptionText(strings.Join(typedElement, ","), false)); err != nil {
				return
			}
		case error:
			if _, _, err = receiver.bot.PostMessage(event.Channel, slack.MsgOptionText(typedElement.Error(), false)); err != nil {
				return
			}
		}
	}
	return nil
}
