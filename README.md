minecraft-bot
===

![License](https://img.shields.io/github/license/ShotaKitazawa/minecraft-bot)
![test](https://github.com/ShotaKitazawa/minecraft-bot/workflows/test/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/ShotaKitazawa/minecraft-bot)
![Dependabot](https://badgen.net/dependabot/ShotaKitazawa/minecraft-bot?icon=dependabot)


minecraft-bot has the following features.

* **notification Bot** : notify Minecraft login/logout events to some chat-tool (LINE/Slack/Discord)
* **operation of Minecraft whitelist** : operate to whitelist of Minecraft Server through ChatBot
* **Prometheus exporter** : support some Minecraft metrics in Prometheus exporter format
* **source is only RCON** : minecraft-bot is not required any Mod, not required reading Minecraft world data, only using RCON

## Demo

### operation by chat & Login/Logout notification to chat

![demo.gif](./images/demo_1.gif)

### using multi chat platform

![demo.gif](./images/demo_2.gif)

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
  -f string
        TOML configuration filepath
  -v    show application version
```

* TOML Configuration File

```
# Minecraft Server Hostname (requirement)
minecraft-hostname = "your_domain"

# basic setting (option)
bind-addr = "0.0.0.0"  # default: "0.0.0.0"
bind-port = 8080       # default: 8080
log-level = "info"     # default: "info",  support "debug", "info", "warn", or "error"


[bot]
# bot basic Configuration (option)
notification-mode = "XXX"  # default: "all", support "none", or "all"

[[bot.line]]
# LINE Bot Configuration (requirement)
endpoint = "/linebot"
channel-secret = "XXX"
channel-token = "XXX"

# LINE Bot Configuration (option)
group-ids = "XXX"  # default: none, cannot push notification without this


[[bot.slack]]
# Slack Bot Configuration (requirement)
token = "XXX"

# Slack Bot Configuration (option)
channel-ids = "XXX"  # default: none, cannot push notification without this


[[bot.discord]]
# Discord Bot Configuration (requirement)
token = "XXX"

# Discord Bot Configuration (option)
channel-ids = "XXX"  # default: none, cannot push notification without this


[rcon]
# connect in RCON to Minecraft (option)
host = "127.0.0.1"    # default: "127.0.0.1"
port = 25575          # default: 25575

# RCON password (requirement)
password = "XXX"


[sharedmem]
# place to store state (support "redis" (recommended), or "local")
mode = "redis"        # default: "local"


[sharedmem.redis]
# Redis info (option if sharedmem.mode == "redis")
host = "127.0.0.1"    # default: "127.0.0.1"
port = 6379           # default: 6379
```

## For more informations

https://github.com/ShotaKitazawa/minecraft-bot/blob/master/docs

## Reference

* [Mod 無し Minecraft で動く Chat Bot + α](https://zenn.dev/kanatakita/articles/5883c5de1b40e17febad)

## Author

[twitter](https://twitter.com/kanatakita)

## Licence

[MIT](https://github.com/ShotaKitazawa/minecraft-bot/blob/master/LICENSE)
