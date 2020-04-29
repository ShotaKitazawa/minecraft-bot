package command

import (
	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/sirupsen/logrus"
)

type PluginHelp struct {
	Logger *logrus.Logger
}

func (p PluginHelp) CommandName() string {
	return `help`
}

func (p PluginHelp) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	queue = append(queue, i18n.T.Sprintf(i18n.MessageHelp))

	return &botplug.MessageOutput{Queue: queue}
}
