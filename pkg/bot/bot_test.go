package bot

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/bot/command"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/mock"
)

const (
	pluginReturnMsgForTest   = `valid`
	MinecraftHostnameForTest = `minecraft`
)

var (
	loggerForTest = logrus.New()
)

func TestBot(t *testing.T) {

	t.Run(`New()`, func(t *testing.T) {
		t.Run(`valid (notificationMode is none)`, func(t *testing.T) {
			notificationMode := "none"
			expected := pcForTest
			expected.NotificationMode = notificationMode

			pc, err := New(loggerForTest, &mock.SharedmemMockValid{}, &mock.RconClientMockValid{}, mock.MockMinecraftHostnameValue, notificationMode)
			assert.Nil(t, err)
			assert.Equal(t, expected, pc)
		})
		t.Run(`valid (notificationMode is all)`, func(t *testing.T) {
			notificationMode := "all"
			expected := pcForTest
			expected.NotificationMode = notificationMode
			expected.Subscriber = &mock.SubscriberMockValid{}

			pc, err := New(loggerForTest, &mock.SharedmemMockValid{}, &mock.RconClientMockValid{}, mock.MockMinecraftHostnameValue, notificationMode)
			assert.Nil(t, err)
			assert.Equal(t, expected, pc)
		})
		t.Run(`invalid (sharedmem is invalid)`, func(t *testing.T) {
			notificationMode := "all"
			expected := pcForTest
			expected.NotificationMode = notificationMode

			_, err := New(loggerForTest, &mock.SharedmemMockInvalid{}, &mock.RconClientMockValid{}, mock.MockMinecraftHostnameValue, notificationMode)
			assert.NotNil(t, err)
		})

	})

	t.Run(`ReceiveMessageEntry()`, func(t *testing.T) {
		t.Run(`succeed`, func(t *testing.T) {
			pc := PluginConfig{
				MinecraftHostname: MinecraftHostnameForTest,
				SharedMem:         &mock.SharedmemMockValid{},
				Plugins:           []PluginInterface{PluginMock{}},
			}
			output := pc.ReceiveMessageEntry(&botplug.MessageInput{
				Messages: []string{`!test`}, // valid
			})
			assert.Equal(t, 1, len(output.Queue))
			result, ok := output.Queue[0].(string)
			assert.True(t, ok)
			assert.Equal(t, pluginReturnMsgForTest, result)
		})
		t.Run(`fail (no such command)`, func(t *testing.T) {
			pc := PluginConfig{
				MinecraftHostname: MinecraftHostnameForTest,
				SharedMem:         &mock.SharedmemMockValid{},
				Plugins:           []PluginInterface{PluginMock{}},
			}
			output := pc.ReceiveMessageEntry(&botplug.MessageInput{
				Messages: []string{`!hoge`}, // no such command
			})
			assert.Equal(t, 1, len(output.Queue))
			result, ok := output.Queue[0].(string)
			assert.True(t, ok)
			assert.Equal(t, i18n.T.Sprintf(i18n.MessageNoSuchCommand), result)
		})
		t.Run(`fail (not command)`, func(t *testing.T) {
			pc := PluginConfig{
				MinecraftHostname: MinecraftHostnameForTest,
				Logger:            loggerForTest,
				Plugins:           []PluginInterface{PluginMock{}},
			}
			output := pc.ReceiveMessageEntry(&botplug.MessageInput{
				Messages: []string{`hoge`}, // not command
			})
			// TODO: dont return nil
			assert.Nil(t, output)
		})
	})
	t.Run(`ReceiveMemberJoinEntry()`, func(t *testing.T) {
		pc := PluginConfig{
			MinecraftHostname: MinecraftHostnameForTest,
			Logger:            loggerForTest,
			Plugins:           []PluginInterface{PluginMock{}},
		}
		output := pc.ReceiveMemberJoinEntry(&botplug.MessageInput{})
		assert.Equal(t, 1, len(output.Queue))
		result, ok := output.Queue[0].(string)
		assert.True(t, ok)
		assert.Equal(t, i18n.T.Sprintf(i18n.MessageMemberJoined, pc.MinecraftHostname), result)
	})

	t.Run(`PushMessageEntry()`, func(t *testing.T) {
		t.Run(`normal (get value from subscriber & push notification)`, func(t *testing.T) {
			pc := PluginConfig{
				MinecraftHostname: MinecraftHostnameForTest,
				Logger:            loggerForTest,
				Sender:            &mock.BotSenderMockValid{},
				Subscriber:        &mock.SubscriberMockValid{},
			}
			queue := pc.PushMessageEntry()
			assert.Equal(t, &botplug.MessageOutput{Queue: []interface{}{mock.MockMessageValue}}, queue)
		})
		t.Run(`normal (subscriber is empty)`, func(t *testing.T) {
			pc := PluginConfig{
				MinecraftHostname: MinecraftHostnameForTest,
				Logger:            loggerForTest,
				Sender:            &mock.BotSenderMockValid{},
			}
			queue := pc.PushMessageEntry()
			assert.Equal(t, &botplug.MessageOutput{}, queue)
		})
		t.Run(`abnormal (subscriber is invalid)`, func(t *testing.T) {
			pc := PluginConfig{
				MinecraftHostname: MinecraftHostnameForTest,
				Logger:            loggerForTest,
				Sender:            &mock.BotSenderMockValid{},
				Subscriber:        &mock.SubscriberMockInvalid{},
			}
			queue := pc.PushMessageEntry()
			assert.Equal(t, &botplug.MessageOutput{}, queue)
		})
	})

	t.Run(`pushToChat()`, func(t *testing.T) {
		pc := PluginConfig{
			MinecraftHostname: MinecraftHostnameForTest,
			Logger:            loggerForTest,
			Sender:            &mock.BotSenderMockValid{},
		}
		queue := pc.pushToChat(`hoge`)
		assert.Equal(t, &botplug.MessageOutput{Queue: []interface{}{"hoge"}}, queue)
	})

}

type PluginMock struct{}

func (p PluginMock) CommandName() string { return `test` }
func (p PluginMock) ReceiveMessage(*botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}
	return &botplug.MessageOutput{Queue: append(queue, pluginReturnMsgForTest)}
}

var pcForTest = &PluginConfig{
	MinecraftHostname: mock.MockMinecraftHostnameValue,
	SharedMem:         &mock.SharedmemMockValid{},
	Rcon:              &mock.RconClientMockValid{},
	Logger:            loggerForTest,
	Plugins: []PluginInterface{
		command.PluginList{
			SharedMem: &mock.SharedmemMockValid{},
			Logger:    loggerForTest,
		},
		command.PluginTitle{
			Rcon:   &mock.RconClientMockValid{},
			Logger: loggerForTest,
		},
		command.PluginWhitelist{
			SharedMem: &mock.SharedmemMockValid{},
			Rcon:      &mock.RconClientMockValid{},
			Logger:    loggerForTest,
		},
		command.PluginHelp{
			Logger: loggerForTest,
		},
		command.PluginID{
			Logger: loggerForTest,
		},
	},
}
