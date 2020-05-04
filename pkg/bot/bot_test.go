package bot

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

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
	t.Run(`ReceiveMessageEntry()`, func(t *testing.T) {
		t.Run(`succeed`, func(t *testing.T) {
			pc := PluginConfig{
				MinecraftHostname: MinecraftHostnameForTest,
				SharedMem:         &mock.SharedmemMockValid{},
				Plugins:           []PluginInterface{PluginMock{}},
			}
			output := pc.ReceiveMessageEntry(&botplug.MessageInput{
				Messages: []string{`/test`}, // valid
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
				Messages: []string{`/hoge`}, // no such command
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
