package command

import (
	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/rcon"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

type PluginWhitelist struct {
	SharedMem sharedmem.SharedMem
	Rcon      rcon.RconClient
	Logger    *logrus.Logger
}

func (p PluginWhitelist) CommandName() string {
	return `whitelist`
}

func (p PluginWhitelist) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	if len(input.Messages) < 2 {
		queue = append(queue, i18n.T.Sprintf(i18n.MessageInvalidArguments))
		return &botplug.MessageOutput{Queue: queue}
	}
	switch input.Messages[1] {
	case `add`:
		if len(input.Messages) < 3 {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageInvalidArguments))
		} else {
			queue = p.add(input.Messages[2:])
		}
	case `delete`:
		if len(input.Messages) < 3 {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageInvalidArguments))
		} else {
			queue = p.delete(input.Messages[2:])
		}
	case `list`:
		queue = p.list()
	default:
		queue = append(queue, i18n.T.Sprintf(i18n.MessageInvalidArguments))
	}

	return &botplug.MessageOutput{Queue: queue}
}

func (p PluginWhitelist) add(users []string) []interface{} {
	var queue []interface{}
	for _, username := range users {
		if p.Rcon.WhitelistAdd(username) != nil {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageUserIncorrect))
		} else {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageWhitelistAdd, username))
		}
	}
	return queue
}

func (p PluginWhitelist) delete(users []string) []interface{} {
	var queue []interface{}
	for _, username := range users {
		if p.Rcon.WhitelistRemove(username) != nil {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageUserIncorrect))
		} else {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageWhitelistRemove, username))
		}
	}
	return queue
}

func (p PluginWhitelist) list() []interface{} {
	var queue []interface{}

	// read data from SharedMem
	data, err := p.SharedMem.SyncReadEntity()
	if err != nil {
		p.Logger.Error(err)
		queue = append(queue, i18n.T.Sprintf(i18n.MessageError))
		return queue
	}

	// whitelist にいるユーザを LINE に送信
	var usernames []string
	for _, username := range data.WhitelistUsernames {
		usernames = append(usernames, username)
	}
	if len(usernames) == 0 {
		queue = append(queue, i18n.T.Sprintf(i18n.MessageNoUserExists))
		return queue
	}
	queue = append(queue, usernames)
	return queue
}
