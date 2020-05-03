package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/minecraft-bot/pkg/bot"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/botplug/line"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/eventer"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/exporter"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/flag"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/rcon"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem/localmem"
	"github.com/ShotaKitazawa/minecraft-bot/pkg/sharedmem/redis"
)

var (
	// These variables are set in build step
	Version  = "unset"
	Revision = "unset"
)

var logger = logrus.New()

func newLogger(loglevel string) *logrus.Logger {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	switch loglevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.Fatal(fmt.Errorf("newLogger: invalid arguments"))
	}
	return logger
}

func main() {
	var err error

	// parse arguments
	conf, err := flag.ArgParse(Version, Revision)
	if err != nil {
		log.Fatal(err)
	}

	// set logger
	logger = newLogger(conf.LogLevel)

	// set LINE config
	var bots []*line.BotAdaptor
	for _, lineConfig := range conf.Bot.LINEConfigs {
		bot, err := line.New(logger, lineConfig.Endpoint, lineConfig.ChannelSecret, lineConfig.ChannelToken, lineConfig.GroupIDs)
		if err != nil {
			logger.Fatal(err)
		}
		bots = append(bots, bot)
	}

	// run sharedMem & get sharedMem instance
	m := func(sharedmemMode string) sharedmem.SharedMem {
		switch sharedmemMode {
		case "local":
			m, err := localmem.New(logger)
			if err != nil {
				logger.Fatal(err)
			}
			return m
		case "redis":
			m, err := redis.New(logger, conf.SharedMem.RedisConfig.Host, conf.SharedMem.RedisConfig.Port)
			if err != nil {
				logger.Fatal(err)
			}
			return m
		default:
			panic(fmt.Errorf("sharedmemMode mismatch"))
		}
	}(conf.SharedMem.Mode)

	// get rcon instance
	rcon, err := rcon.New(conf.Rcon.Host, conf.Rcon.Port, conf.Rcon.Password)
	if err != nil {
		logger.Fatal(err)
	}

	// run eventer
	eventer, err := eventer.New(conf.MinecraftHostname, m, rcon, logger)
	if err != nil {
		logger.Fatal(err)
	}
	go eventer.Run()

	// run exporter
	collector, err := exporter.New(m, logger)
	if err != nil {
		logger.Fatal(err)
	}
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())

	// run bot
	for _, botInstance := range bots {
		handler, err := botInstance.WithPlugin(
			bot.New(conf.MinecraftHostname, m, rcon, logger),
		).NewHandler()
		if err != nil {
			logger.Fatal(err)
		}
		http.Handle(botInstance.Endpoint, handler)
	}

	logger.Fatal(http.ListenAndServe(":8080", nil))
}
