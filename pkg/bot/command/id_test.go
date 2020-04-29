package command

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
)

const (
	userIDForTestCommandID  = `user`
	groupIDForTestCommandID = `group`
)

var (
	loggerForTestCommandID = logrus.New()
)

func NewPluginIDForTest() PluginID {
	return PluginID{Logger: loggerForTestCommandID}
}

func TestCommandID(t *testing.T) {
	t.Run(`ReceiveMessage()`, func(t *testing.T) {
		p := NewPluginIDForTest()
		output := p.ReceiveMessage(&botplug.MessageInput{
			Source: &botplug.Source{
				UserID:  userIDForTestCommandID,
				GroupID: groupIDForTestCommandID,
			},
		})
		assert.Equal(t, 1, len(output.Queue))
		result, ok := output.Queue[0].(string)
		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf(messageID, userIDForTestCommandID, groupIDForTestCommandID), result)
	})
}
