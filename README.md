minecraft-bot
===

![License](https://img.shields.io/github/license/ShotaKitazawa/minecraft-bot)
![test](https://github.com/ShotaKitazawa/minecraft-bot/workflows/test/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/ShotaKitazawa/minecraft-bot)
![Dependabot](https://flat.badgen.net/dependabot/thepracticaldev/dev.to?icon=dependabot)
![Codecov](https://codecov.io/gh/ShotaKitazawa/minecraft-bot/branch/master/graph/badge.svg)


minecraft-bot has the following features.

* **notification Bot** : notify Minecraft login/logout events to some chat-tool (LINE/Slack/Discord)
* **Prometheus exporter** : support some Minecraft metrics in Prometheus exporter format
* **source is only RCON** : minecraft-bot is not required any Mod, not required reading Minecraft world data, only using RCON

## Demo

![demo.gif](./images/demo.gif)

### Dashboard

* using `./dashboard/minecraft_users.json`
    * pre-install `vonage-status-panel` Panel Plugin

![minecraft_users](./images/minecraft_users.png)

## Installation

* only download binary from [release](https://github.com/ShotaKitazawa/minecraft-bot/releases)

* or execute bellow command

```
go get -u github.com/ShotaKitazawa/minecraft-bot/cmd/minecraft-bot
```

## Usage

```
$ minecraft-bot -h
Usage of minecraft-bot:
  -line-channel-secret="": LINE Bot's Channel Secret
  -line-channel-token="": LINE Bot's Channel Token
  -line-group-id="": specified LINE Group ID, send push message to this Group
  -log-level="info": Log Level (debug, info, warn, error)
  -minecraft-hostname="": Minecraft Hostname
  -rcon-host="": RCON Host
  -rcon-password="": RCON Password
  -rcon-port=25575: RCON Port
  -redis-host="127.0.0.1": Redis Host (enabled when sharedmem-mode=redis)
  -redis-port=6379: Redis Port (enabled when sharedmem-mode=redis)
  -sharedmem-mode="local": using Shared Memory ("local" or "redis")
  -v=false: show application version
```

### Execution example


* enable LINE Bot
    * run Bot on the same server as Minecraft
* using Redis in sharedmem
    * run Bot on the same server as Redis

```
$ minecraft-bot -line-channel-secret=XXX -line-channel-token=XXX -line-group-id=XXX -rcon-host=127.0.0.1 -rcon-password=XXX -redis-host=127.0.0.1 -sharedmem-mode=redis
```


## Setup

### Minecraft Server Setup

In your Minecraft `server.properties`, make sure you have and restart the server.

```
enable-rcon=true
rcon.password=[minecraftRconPassword]
rcon.port=[minecraftRconPort]
```

### Bot do not support HTTPS

This bot run HTTP server, but Webhook configuration required HTTPS in most chat-provider.
Please following the below.

* using HTTPS reverse-proxy server (nginx, Caddy, etc..) & run Bot beside Minecraft server
* using PaaS (Heroku, Google App Engine, etc..)
    * not recommended (RCON connection is not crypted)

### Bot for LINE

* Setup LINE `Messageing API` : https://developers.line.biz/console/
    * Webhook URL: `https://<your_domain>/linebot`

### Bot for Slack

TBD

### Bot for Discord

TBD

## Architecture

![](./images/architecture.png)

