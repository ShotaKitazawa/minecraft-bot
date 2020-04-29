package line

import (
	"strings"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/sirupsen/logrus"
)

type Config struct {
	botplug.Config
	ChannelSecret string
	ChannelToken  string
	bot           *linebot.Client
}

func New(logger *logrus.Logger, channelSecret, channelToken string, groupIDsStr string) (*BotReceiver, *BotSender, error) {
	botplugLINEConfig := botplug.New(logger)
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil, nil, err
	}
	botConfig := Config{
		Config:        botplugLINEConfig,
		ChannelSecret: channelSecret,
		ChannelToken:  channelToken,
		bot:           bot,
	}
	return &BotReceiver{
			Config: botConfig,
		}, &BotSender{
			Config:   botConfig,
			GroupIDs: strings.Split(groupIDsStr, ","),
		}, nil
}

func connect(config Config) (*linebot.Client, error) {
	bot, err := linebot.New(
		config.ChannelSecret,
		config.ChannelToken,
	)
	if err != nil {
		return nil, err
	}
	return bot, nil

}
