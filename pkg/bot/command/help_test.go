package command

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/domain/i18n"
)

var (
	loggerForTestCommandHelp = logrus.New()
)

func NewPluginHelpForTest() PluginHelp {
	return PluginHelp{Logger: loggerForTestCommandHelp}
}

func TestCommandHelp(t *testing.T) {
	t.Run(`ReceiveMessage()`, func(t *testing.T) {
		p := NewPluginHelpForTest()
		output := p.ReceiveMessage(&botplug.MessageInput{})
		assert.Equal(t, 1, len(output.Queue))
		result, ok := output.Queue[0].(string)
		assert.True(t, ok)
		assert.Equal(t, i18n.T.Sprintf(i18n.MessageHelp), result)
	})
}
