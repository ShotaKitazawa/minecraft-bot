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
	loggerForTestCommandWhitelist = logrus.New()
	userForTestCommandWhitelist   = domain.User{
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

func TestCommandWhitelist(t *testing.T) {

	t.Run(`ReceiveMessage()`, func(t *testing.T) {

		t.Run(`normal`, func(t *testing.T) {

			t.Run(`subcommand "add"`, func(t *testing.T) {

				t.Run(`succeed`, func(t *testing.T) {
					p := PluginWhitelist{
						Logger: loggerForTestCommandWhitelist,
						Rcon: &mock.RconClientMockValid{
							LoginUsernames: []string{mock.MockUserNameValue},
						},
						SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
							WhitelistUsernames: []string{mock.MockUserNameValue},
						}},
					}
					output := p.ReceiveMessage(&botplug.MessageInput{
						Messages: []string{`whitelist`, `add`, mock.MockUserNameValue}, // valid (len == 3)
					})
					assert.Equal(t, 1, len(output.Queue))
					result, ok := output.Queue[0].(string)
					assert.True(t, ok)
					assert.Equal(t, i18n.T.Sprintf(i18n.MessageWhitelistAdd, mock.MockUserNameValue), result)
				})
				t.Run(`fail (no user exist)`, func(t *testing.T) {
					// TBD
				})
				t.Run(`fail (invalid arguments)`, func(t *testing.T) {
					p := PluginWhitelist{
						Logger: loggerForTestCommandWhitelist,
						Rcon: &mock.RconClientMockValid{
							LoginUsernames: []string{mock.MockUserNameValue},
						},
						SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
							WhitelistUsernames: []string{mock.MockUserNameValue},
						}},
					}
					output := p.ReceiveMessage(&botplug.MessageInput{
						Messages: []string{`whitelist`, `add`}, // invalid (len < 3)
					})
					assert.Equal(t, 1, len(output.Queue))
					result, ok := output.Queue[0].(string)
					assert.True(t, ok)
					assert.Equal(t, i18n.T.Sprintf(i18n.MessageInvalidArguments), result)
				})

			})

			t.Run(`subcommand "delete"`, func(t *testing.T) {

				t.Run(`succeed`, func(t *testing.T) {
					p := PluginWhitelist{
						Logger: loggerForTestCommandWhitelist,
						Rcon: &mock.RconClientMockValid{
							WhitelistedUsernames: []string{mock.MockUserNameValue},
						},
						SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
							WhitelistUsernames: []string{mock.MockUserNameValue},
						}},
					}
					output := p.ReceiveMessage(&botplug.MessageInput{
						Messages: []string{`whitelist`, `delete`, mock.MockUserNameValue}, // valid (len == 3)
					})
					assert.Equal(t, 1, len(output.Queue))
					result, ok := output.Queue[0].(string)
					assert.True(t, ok)
					assert.Equal(t, i18n.T.Sprintf(i18n.MessageWhitelistRemove, mock.MockUserNameValue), result)
				})
				t.Run(`fail (no user exist)`, func(t *testing.T) {
					p := PluginWhitelist{
						Logger: loggerForTestCommandWhitelist,
						Rcon:   &mock.RconClientMockValid{},
						SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
							WhitelistUsernames: []string{mock.MockUserNameValue},
						}},
					}
					output := p.ReceiveMessage(&botplug.MessageInput{
						Messages: []string{`whitelist`, `delete`, mock.MockUserNameValue}, // valid (len == 3)
					})
					assert.Equal(t, 1, len(output.Queue))
					result, ok := output.Queue[0].(string)
					assert.True(t, ok)
					assert.Equal(t, i18n.T.Sprintf(i18n.MessageUserIncorrect), result)
				})
				t.Run(`fail (invalid arguments)`, func(t *testing.T) {
					p := PluginWhitelist{
						Logger: loggerForTestCommandWhitelist,
						Rcon: &mock.RconClientMockValid{
							LoginUsernames: []string{mock.MockUserNameValue},
						},
						SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
							WhitelistUsernames: []string{mock.MockUserNameValue},
						}},
					}
					output := p.ReceiveMessage(&botplug.MessageInput{
						Messages: []string{`whitelist`, `delete`}, // valid (len < 3)
					})
					assert.Equal(t, 1, len(output.Queue))
					result, ok := output.Queue[0].(string)
					assert.True(t, ok)
					assert.Equal(t, i18n.T.Sprintf(i18n.MessageInvalidArguments), result)
				})
			})

			t.Run(`subcommand "list"`, func(t *testing.T) {
				t.Run(`succeed`, func(t *testing.T) {
					p := PluginWhitelist{
						Logger: loggerForTestCommandWhitelist,
						Rcon: &mock.RconClientMockValid{
							LoginUsernames: []string{mock.MockUserNameValue},
						},
						SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
							WhitelistUsernames: []string{mock.MockUserNameValue},
						}},
					}
					output := p.ReceiveMessage(&botplug.MessageInput{
						Messages: []string{`whitelist`, `list`},
					})
					assert.Equal(t, 1, len(output.Queue))
					result, ok := output.Queue[0].([]string)
					assert.True(t, ok)
					assert.Equal(t, []string{mock.MockUserNameValue}, result)
				})
				t.Run(`fail (no user exist)`, func(t *testing.T) {
					p := PluginWhitelist{
						Logger:    loggerForTestCommandWhitelist,
						Rcon:      &mock.RconClientMockValid{},
						SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{}},
					}
					output := p.ReceiveMessage(&botplug.MessageInput{
						Messages: []string{`whitelist`, `list`},
					})
					assert.Equal(t, 1, len(output.Queue))
					result, ok := output.Queue[0].(string)
					assert.True(t, ok)
					assert.Equal(t, i18n.T.Sprintf(i18n.MessageNoUserExists), result)
				})
			})

			t.Run(`invalid subcommand`, func(t *testing.T) {
			})

		})

		t.Run(`abnormal (input)`, func(t *testing.T) {
			p := PluginWhitelist{
				Logger: loggerForTestCommandWhitelist,
				Rcon: &mock.RconClientMockValid{
					LoginUsernames: []string{mock.MockUserNameValue},
				},
				SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
					WhitelistUsernames: []string{mock.MockUserNameValue},
				}},
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
			t.Run(`subcommand "add"`, func(t *testing.T) {
				p := PluginWhitelist{
					Logger: loggerForTestCommandWhitelist,
					Rcon:   &mock.RconClientMockInvalid{}, // abnormal
					SharedMem: &mock.SharedmemMockValid{Data: &domain.Entity{
						WhitelistUsernames: []string{mock.MockUserNameValue},
					}},
				}

				output := p.ReceiveMessage(&botplug.MessageInput{
					Messages: []string{`whitelist`, `add`, mock.MockUserNameValue},
				})
				assert.Equal(t, 1, len(output.Queue))
				result, ok := output.Queue[0].(string)
				assert.True(t, ok)
				assert.Equal(t, i18n.T.Sprintf(i18n.MessageUserIncorrect), result)
			})
		})
		t.Run(`abnormal (sharedmem)`, func(t *testing.T) {
			t.Run(`subcommand "list"`, func(t *testing.T) {
				p := PluginWhitelist{
					Logger: loggerForTestCommandWhitelist,
					Rcon: &mock.RconClientMockValid{
						LoginUsernames: []string{mock.MockUserNameValue},
					},
					SharedMem: &mock.SharedmemMockInvalid{},
				}

				output := p.ReceiveMessage(&botplug.MessageInput{
					Messages: []string{`whitelist`, `list`},
				})
				assert.Equal(t, 1, len(output.Queue))
				result, ok := output.Queue[0].(string)
				assert.True(t, ok)
				assert.Equal(t, i18n.T.Sprintf(i18n.MessageError), result)
			})
		})

	})
}
