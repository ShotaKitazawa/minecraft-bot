package bot

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/bot/command"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/rcon"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
)

const (
	commandPrefix = `/`
)

type PluginConfig struct {
	MinecraftHostname string
	SharedMem         sharedmem.SharedMem
	Subscriber        sharedmem.Subscriber
	Rcon              rcon.RconClient
	Logger            *logrus.Logger
	Plugins           []PluginInterface
	Sender            botplug.BotSender
}

func New(minecraftHostname string, m sharedmem.SharedMem, rcon *rcon.Client, logger *logrus.Logger) (*PluginConfig, error) {
	subscriber, err := m.NewSubscriber()
	if err != nil {
		return nil, err
	}
	pc := &PluginConfig{
		MinecraftHostname: minecraftHostname,
		SharedMem:         m,
		Subscriber:        subscriber,
		Rcon:              rcon,
		Logger:            logger,
		Plugins: []PluginInterface{
			command.PluginList{
				SharedMem: m,
				Logger:    logger,
			},
			command.PluginTitle{
				Rcon:   rcon,
				Logger: logger,
			},
			command.PluginWhitelist{
				SharedMem: m,
				Rcon:      rcon,
				Logger:    logger,
			},
			command.PluginHelp{
				Logger: logger,
			},
			command.PluginID{
				Logger: logger,
			},
		},
	}
	return pc, nil
}

func (pc *PluginConfig) ReceiveMessageEntry(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	if !strings.HasPrefix(input.Messages[0], commandPrefix) {
		// TODO: dont return nil
		return nil
	}
	input.Messages[0] = strings.TrimLeft(input.Messages[0], commandPrefix)

	if input.Source != nil {
		pc.Logger.WithFields(logrus.Fields{
			"source": *input.Source,
		}).Debug(input.Messages)
	}

	for _, plugin := range pc.Plugins {
		if input.Messages[0] == plugin.CommandName() {
			return plugin.ReceiveMessage(input)
		}
	}

	queue = append(queue, i18n.T.Sprintf(i18n.MessageNoSuchCommand))
	return &botplug.MessageOutput{Queue: queue}
}

func (pc *PluginConfig) ReceiveMemberJoinEntry(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}
	queue = append(queue, i18n.T.Sprintf(i18n.MessageMemberJoined, pc.MinecraftHostname))
	return &botplug.MessageOutput{Queue: queue}
}

func (pc *PluginConfig) PushMessageEntry() *botplug.MessageOutput {
	message, err := pc.Subscriber.SyncSubscribeMessage() // wait until get data
	if err != nil {
		pc.Logger.Error(err)
	}
	queue := pc.pushToChat(message.Msg)
	if err != nil {
		pc.Logger.Error(err)
	}
	return queue
}

func (pc *PluginConfig) pushToChat(msg string) *botplug.MessageOutput {
	var queue []interface{}

	queue = append(queue, msg)
	return &botplug.MessageOutput{Queue: queue}
}
