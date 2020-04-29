package botplug

import "github.com/sirupsen/logrus"

type Config struct {
	ID     uint
	Logger *logrus.Logger
}

var id = 0

func New(logger *logrus.Logger) Config {
	id++
	return Config{
		ID:     uint(id),
		Logger: logger,
	}
}
