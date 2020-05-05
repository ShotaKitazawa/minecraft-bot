package botplug

import (
	"github.com/sirupsen/logrus"
)

type BaseConfig struct {
	ID      uint
	Logger  *logrus.Logger
	Plugins []BotPlugin
}

var id = 0

func New(logger *logrus.Logger) BaseConfig {
	id++
	return BaseConfig{
		ID:     uint(id),
		Logger: logger,
	}
}
