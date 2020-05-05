package slack

import (
	"strconv"
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
func (sender *BotAdaptor) SendTextMessageToChannels(text string) (err error) {
	for _, channelID := range sender.ChannelIDs {
		if err := sender.sendTextMessage(channelID, text); err != nil {
			return err
		}
	}
	return nil
}
func (ba *BotAdaptor) sendTextMessage(channel, text string) (err error) {
	_, _, err = ba.bot.PostMessage(channel, slack.MsgOptionText(text, false))
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
func (ba *BotAdaptor) Run() error {
	rtm := ba.bot.NewRTM()
	go rtm.ManageConnection()
	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := ba.receiveTextMessage(ev); err != nil {
				ba.Logger.Error(err)
			}
		case *slack.ChannelJoinedEvent:
			if err := ba.receiveMemberJoin(ev); err != nil {
				ba.Logger.Error(err)
			}

		}
	}
	return nil
}

// receiveTextMessage is called when Receive Chat Message
func (receiver *BotAdaptor) receiveTextMessage(event *slack.MessageEvent) (err error) {
	unixTime, err := strconv.Atoi(strings.Split(event.Timestamp, ".")[0])
	if err != nil {
		return err
	}
	input := &botplug.MessageInput{
		Timestamp: time.Unix(int64(unixTime), 0),
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
		texts, err := botplug.FormatToText(output)
		if err != nil {
			return err
		}
		for _, text := range texts {
			receiver.sendTextMessage(input.Source.GroupID, text)
		}
	}

	return nil
}

// receiveMemberJoin is called when someone join at Chat Group
func (receiver *BotAdaptor) receiveMemberJoin(event *slack.ChannelJoinedEvent) (err error) {
	input := &botplug.MessageInput{
		//Timestamp: none
		Source: &botplug.Source{
			Type: string(event.Type),
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
			receiver.sendTextMessage(input.Source.GroupID, text)
		}
	}
	return nil
}
