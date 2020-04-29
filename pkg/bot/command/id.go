package command

import (
	"fmt"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/sirupsen/logrus"
)

const (
	messageID = `
UserID: %s
GroupID: %s
`
)

type PluginID struct {
	Logger *logrus.Logger
}

func (p PluginID) CommandName() string {
	return `id`
}

func (p PluginID) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	queue = append(queue, fmt.Sprintf(messageID, input.Source.UserID, input.Source.GroupID))

	return &botplug.MessageOutput{Queue: queue}
}
