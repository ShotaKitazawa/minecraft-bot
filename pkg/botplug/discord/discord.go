package discord

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
)

type BotAdaptor struct {
	botplug.BaseConfig
	Token            string
	NotificationMode string
	ChannelIDs       []string
	bot              *discordgo.Session
}

func New(logger *logrus.Logger, token, channelIDsStr, notificationMode string) (*BotAdaptor, error) {
	botplugSlackConfig := botplug.New(logger)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	ba := &BotAdaptor{
		BaseConfig:       botplugSlackConfig,
		Token:            token,
		NotificationMode: notificationMode,
		ChannelIDs:       strings.Split(channelIDsStr, ","),
		bot:              dg,
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
	_, err = ba.bot.ChannelMessageSend(channel, text)
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
	ba.bot.AddHandler(func(s *discordgo.Session, ev *discordgo.MessageCreate) {
		if ev.Author.ID == s.State.User.ID {
			return
		}
		switch ev.Type {
		case discordgo.MessageTypeDefault:
			if err := ba.receiveTextMessage(ev); err != nil {
				ba.Logger.Error(err)
			}
		case discordgo.MessageTypeGuildMemberJoin:
			if err := ba.receiveMemberJoin(ev); err != nil {
				ba.Logger.Error(err)
			}
		}
	})
	err := ba.bot.Open()
	if err != nil {
		return err
	}

	return nil
}

// receiveTextMessage is called when Receive Chat Message
func (receiver *BotAdaptor) receiveTextMessage(event *discordgo.MessageCreate) (err error) {
	t, err := time.Parse(time.RFC3339Nano, string(event.Timestamp))
	if err != nil {
		return err
	}
	input := &botplug.MessageInput{
		Timestamp: t,
		Source: &botplug.Source{
			Type:    string(event.Type),
			UserID:  event.Author.ID,
			GroupID: event.ChannelID,
		},
		Messages: strings.Fields(event.Content),
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
func (receiver *BotAdaptor) receiveMemberJoin(event *discordgo.MessageCreate) (err error) {
	t, err := time.Parse(time.RFC3339Nano, string(event.Timestamp))
	if err != nil {
		return err
	}
	input := &botplug.MessageInput{
		Timestamp: t,
		Source: &botplug.Source{
			Type:    string(event.Type),
			UserID:  event.Author.ID,
			GroupID: event.ChannelID,
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
