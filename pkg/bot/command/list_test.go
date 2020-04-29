package command

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/mock"
)

var (
	loggerForTestCommandList = logrus.New()
	userForTestCommandList   = domain.User{
		Name:    mock.MockUserNameValue,
		Health:  mock.MockUserHealthValue,
		XpLevel: mock.MockUserXpValue,
		Position: domain.Position{
			X: mock.MockUserPosXValue,
			Y: mock.MockUserPosYValue,
			Z: mock.MockUserPosZValue,
		},
	}
)

func TestCommandList(t *testing.T) {

	t.Run(`ReceiveMessage()`, func(t *testing.T) {

		t.Run(`normal`, func(t *testing.T) {

			t.Run(`1 user login`, func(t *testing.T) {
				p := PluginList{
					Logger: loggerForTestCommandList,
					SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
						AllUsers:   []domain.User{userForTestCommandList},
						LoginUsers: []domain.User{userForTestCommandList},
					}},
				}
				output := p.ReceiveMessage(&botplug.MessageInput{})
				assert.Equal(t, 1, len(output.Queue))
				result, ok := output.Queue[0].([]string)
				assert.True(t, ok)
				assert.Equal(t, []string{mock.MockUserNameValue}, result)
			})

			t.Run(`no users login`, func(t *testing.T) {
				p := PluginList{
					Logger: loggerForTestCommandList,
					SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
						AllUsers: []domain.User{userForTestCommandList},
					}},
				}
				output := p.ReceiveMessage(&botplug.MessageInput{})
				assert.Equal(t, 1, len(output.Queue))
				result, ok := output.Queue[0].(string)
				assert.True(t, ok)
				assert.Equal(t, i18n.T.Sprintf(i18n.MessageNoUserExists), result)
			})

		})

		t.Run(`abnormal (SharedMem)`, func(t *testing.T) {
			p := PluginList{
				Logger:    loggerForTestCommandList,
				SharedMem: &mock.SharedmemMockInvalid{},
			}
			output := p.ReceiveMessage(&botplug.MessageInput{})
			assert.Equal(t, 1, len(output.Queue))
			result, ok := output.Queue[0].(string)
			assert.True(t, ok)
			assert.Equal(t, i18n.T.Sprintf(i18n.MessageError), result)
		})

	})
}
