package bot

import "github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"

type PluginInterface interface {
	CommandName() string
	ReceiveMessage(*botplug.MessageInput) *botplug.MessageOutput
}
