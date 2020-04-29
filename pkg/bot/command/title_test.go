package command

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/mock"
)

var (
	loggerForTestCommandTitle = logrus.New()
)

func TestCommandTitle(t *testing.T) {

	t.Run(`ReceiveMessage()`, func(t *testing.T) {

		t.Run(`normal`, func(t *testing.T) {

			t.Run(`1 user login`, func(t *testing.T) {
				p := PluginTitle{
					Logger: loggerForTestCommandTitle,
					Rcon: &mock.RconClientMockValid{
						LoginUsernames: []string{mock.MockUserNameValue},
					},
				}
				output := p.ReceiveMessage(&botplug.MessageInput{
					Messages: []string{`title`, `test`},
				})
				assert.Equal(t, 1, len(output.Queue))
				result, ok := output.Queue[0].(string)
				assert.True(t, ok)
				assert.Equal(t, i18n.T.Sprintf(i18n.MessageSentMessage, mock.MockUserNameValue), result)
			})

			t.Run(`no users login`, func(t *testing.T) {
				p := PluginTitle{
					Logger: loggerForTestCommandList,
					Rcon:   &mock.RconClientMockValid{},
				}
				output := p.ReceiveMessage(&botplug.MessageInput{
					Messages: []string{`title`, `test`},
				})
				assert.Equal(t, 1, len(output.Queue))
				result, ok := output.Queue[0].(string)
				assert.True(t, ok)
				assert.Equal(t, i18n.T.Sprintf(i18n.MessageNoLoginUserExists), result)
			})

		})

		t.Run(`abnormal (input)`, func(t *testing.T) {
			p := PluginTitle{
				Logger: loggerForTestCommandList,
				Rcon:   &mock.RconClientMockValid{},
			}
			output := p.ReceiveMessage(&botplug.MessageInput{
				Messages: []string{}, // abnormal (len < 2)
			})
			assert.Equal(t, 1, len(output.Queue))
			result, ok := output.Queue[0].(string)
			assert.True(t, ok)
			assert.Equal(t, i18n.T.Sprintf(i18n.MessageInvalidArguments), result)
		})

		t.Run(`abnormal (rcon)`, func(t *testing.T) {
			p := PluginTitle{
				Logger: loggerForTestCommandList,
				Rcon:   &mock.RconClientMockInvalid{}, // abnormal
			}
			output := p.ReceiveMessage(&botplug.MessageInput{
				Messages: []string{`title`, `test`},
			})
			assert.Equal(t, 1, len(output.Queue))
			result, ok := output.Queue[0].(string)
			assert.True(t, ok)
			assert.Equal(t, i18n.T.Sprintf(i18n.MessageError), result)
		})

	})
}
