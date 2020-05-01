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
		panic(fmt.Errorf("newLogger: invalid arguments"))
	}
	return logger
}

func main() {
	var err error

	// parse arguments
	conf, err := flag.ArgParse(Version, Revision)
	// set logger
	logger = newLogger(conf.LogLevel)

	// error check of "parse arguments"
	if err != nil {
		logger.Fatal(err)
	}

	// set LINE config
	botReceiver, botSender, err := line.New(logger, conf.Bot.LINEConfig.ChannelSecret, conf.Bot.LINEConfig.ChannelToken, conf.Bot.LINEConfig.GroupIDs)

	// run sharedMem & get sharedMem instance
	m := func(sharedmemMode string) sharedmem.SharedMem {
		switch sharedmemMode {
		case "local":
			m, err := localmem.New(logger)
			if err != nil {
				panic(err)
			}
			return m
		case "redis":
			m, err := redis.New(logger, conf.SharedMem.RedisConfig.Host, conf.SharedMem.RedisConfig.Port)
			if err != nil {
				panic(err)
			}
			return m
		default:
			panic(fmt.Errorf("sharedmemMode mismatch"))
		}
	}(conf.SharedMem.Mode)

	// get rcon instance
	rcon, err := rcon.New(conf.Rcon.Host, conf.Rcon.Port, conf.Rcon.Password)
	if err != nil {
		panic(err)
	}

	// run eventer
	eventer, err := eventer.New(conf.MinecraftHostname, botSender, m, rcon, logger)
	if err != nil {
		panic(err)
	}
	go eventer.Run()

	// run exporter
	collector, err := exporter.New(m, logger)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())

	// run bot
	handler, err := botReceiver.WithPlugin(
		bot.New(conf.MinecraftHostname, m, rcon, logger),
	).NewHandler()
	if err != nil {
		panic(err)
	}
	http.Handle("/linebot", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
