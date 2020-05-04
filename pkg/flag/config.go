package flag

import (
	"errors"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BindAddr          string          `toml:"bind-addr"`
	BindPort          int             `toml:"bind-port"`
	MinecraftHostname string          `toml:"minecraft-hostname"`
	LogLevel          string          `toml:"log-level"`
	Bot               BotConfig       `toml:"bot"`
	Rcon              RconConfig      `toml:"rcon"`
	SharedMem         SharedMemConfig `toml:"sharedmem"`
}

type BotConfig struct {
	LINEConfigs []LINEConfig `toml:"line"`
	// TBD
	//SlackConfigs        []SlackConfig   `toml:"slack"`
	//DiscordConfigs      []DiscordConfig `toml:"discord"`
}

type LINEConfig struct {
	Endpoint         string `toml:"endpoint"`
	ChannelSecret    string `toml:"channel-secret"`
	ChannelToken     string `toml:"channel-token"`
	GroupIDs         string `toml:"group-ids"`
	NotificationMode string `toml:"notification-mode"`
}

// TBD
// type SlackConfig struct{}
// type DiscordConfig struct{}

type RconConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
}

type SharedMemConfig struct {
	Mode        string      `toml:"mode"`
	RedisConfig RedisConfig `toml:"redis"`
}

type RedisConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func ParseConfig(filepath string) (*Config, error) {
	handle, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	conf, err := ScanConfig(handle)
	if err != nil {
		return nil, err
	}
	return conf, ValidateConfig(conf)
}

func ScanConfig(r io.Reader) (*Config, error) {
	var config *Config
	_, err := toml.DecodeReader(r, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func ValidateConfig(config *Config) error {
	if config.BindAddr == "" {
		config.BindAddr = "0.0.0.0"
	}
	if config.BindPort == 0 {
		config.BindPort = 8080
	}
	if config.MinecraftHostname == "" {
		return errors.New(`"minecraft-hostname" is requirement field`)
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	} else if !(config.LogLevel == "debug" ||
		config.LogLevel == "info" ||
		config.LogLevel == "warn" ||
		config.LogLevel == "error") {
		config.LogLevel = "debug"
		return errors.New(`"log-level" only support "debug", "info", "warn", and "error"`)
	}

	for _, LINEConfig := range config.Bot.LINEConfigs {
		if LINEConfig.Endpoint == "" {
			return errors.New(`"bot.line.endpoint" is requirement field`)
		}
		if LINEConfig.ChannelSecret == "" {
			return errors.New(`"bot.line.channel-secret" is requirement field`)
		}
		if LINEConfig.ChannelToken == "" {
			return errors.New(`"bot.line.channel-token" is requirement field`)
		}
		if LINEConfig.GroupIDs == "" {
			logger.Warnf(`"bot.line.group-id" is empty, push notification is disabled.`)
		}
		if LINEConfig.NotificationMode == "" {
			LINEConfig.NotificationMode = "all"
		} else if !(LINEConfig.NotificationMode == "all" ||
			LINEConfig.NotificationMode == "none") {
			return errors.New(`"notification-mode" only support "all", and "none"`)
		}
	}
	if config.Rcon.Host == "" {
		config.Rcon.Host = "127.0.0.1"
	}
	if config.Rcon.Port == 0 {
		config.Rcon.Port = 25575
	}
	if config.Rcon.Password == "" {
		return errors.New(`"rcon.password" is requirement field`)
	}
	if config.SharedMem.Mode == "" {
		config.SharedMem.Mode = "local"
	} else if !(config.SharedMem.Mode == "local" || config.SharedMem.Mode == "redis") {
		return errors.New(`"sharedmem.mode" only support "local", and "redis"`)
	}
	if config.SharedMem.Mode == `redis` {
		if config.SharedMem.RedisConfig.Host == "" {
			config.SharedMem.RedisConfig.Host = "127.0.0.1"
		}
		if config.SharedMem.RedisConfig.Port == 0 {
			config.SharedMem.RedisConfig.Port = 6379
		}
	}

	return nil
}
