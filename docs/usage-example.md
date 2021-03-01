## Usage example

* enable one LINE Bot belong to 2 Groups
* using Redis in sharedmem
* minecraft-bot, Minecraft, Redis exist in the same server

```
$ cat config.toml
minecraft-hostname = "your_domain"

[[bot.line]]
endpoint = "/linebot"
channel-secret = "XXX"
channel-token = "XXX"
group-ids = "GROUP1,GROUP2"

[rcon]
password = "XXX"

[sharedmem]
mode = "redis"
```

* enable two LINE Bot & one Slack Bot, each belong to 1 group
* using Redis in sharedmem
* minecraft-bot, Minecraft, Redis exist in the same server

```
$ cat config.toml
minecraft-hostname = "your_domain"

[[bot.line]]
endpoint = "/linebot"
channel-secret = "XXX"
channel-token = "XXX"
group-ids = "XXX"

[[bot.line]]
endpoint = "/test"
channel-secret = "XXX"
channel-token = "XXX"
group-ids = "XXX"

[[bot.slack]]
token = "XXX"
channel-ids = "XXX"

[rcon]
password = "XXX"

[sharedmem]
mode = "redis"
```

