package command

import (
	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

type PluginList struct {
	SharedMem sharedmem.SharedMem
	Logger    *logrus.Logger
}

func (p PluginList) CommandName() string {
	return `list`
}

func (p PluginList) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	// read data from SharedMem
	data, err := p.SharedMem.SyncReadEntityFromSharedMem()
	if err != nil {
		p.Logger.Error(err)
		queue = append(queue, i18n.T.Sprintf(i18n.MessageError))
		return &botplug.MessageOutput{Queue: queue}
	}

	// ログイン中のユーザを LINE に送信
	var loginUsernames []string
	for _, user := range data.LoginUsers {
		loginUsernames = append(loginUsernames, user.Name)
	}
	if loginUsernames == nil {
		queue = append(queue, i18n.T.Sprintf(i18n.MessageNoUserExists))
		return &botplug.MessageOutput{Queue: queue}
	}
	queue = append(queue, loginUsernames)
	return &botplug.MessageOutput{Queue: queue}
}
