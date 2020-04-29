package command

import (
	"strings"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/rcon"
	"github.com/sirupsen/logrus"
)

type PluginTitle struct {
	Rcon   rcon.RconClient
	Logger *logrus.Logger
}

func (p PluginTitle) CommandName() string {
	return `title`
}

func (p PluginTitle) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	if len(input.Messages) < 2 {
		queue = append(queue, i18n.T.Sprintf(i18n.MessageInvalidArguments))
		return &botplug.MessageOutput{Queue: queue}
	}

	// send RCON
	destUsers, err := p.Rcon.Title(strings.Join(input.Messages[1:], " "))
	if err != nil {
		p.Logger.Error(err)
		queue = append(queue, i18n.T.Sprintf(i18n.MessageError))
		return &botplug.MessageOutput{Queue: queue}
	}
	if len(destUsers) == 0 {
		queue = append(queue, i18n.T.Sprintf(i18n.MessageNoLoginUserExists))
		return &botplug.MessageOutput{Queue: queue}
	}
	for _, user := range destUsers {
		queue = append(queue, i18n.T.Sprintf(i18n.MessageSentMessage, user))
	}

	return &botplug.MessageOutput{Queue: queue}
}
