## Required pre-setup

### Minecraft Server Setup

In your Minecraft `server.properties`, make sure you have and restart the server.

```
enable-rcon=true
rcon.password=[minecraftRconPassword]
rcon.port=[minecraftRconPort]
```

### setup Bot for LINE

* Setup LINE `Messageing API` : https://developers.line.biz/console/
    * Webhook URL: `https://<your_domain>/<bot.line.endpoint>`

* Look up GroupID
    1. run bot without `group-ids` of `[[bot.line]]` in config.toml
    2. chat `!id` in specified group
    3. check response of `GroupID: XXX`
    4. postscript `group-ids` of `[[bot.line]]` in config.toml & re-run bot

### setup Bot for Slack

* Setup Slack Bot & Get `Bot User OAuth Access Token` : https://api.slack.com/apps/

* Look up GroupID
    1. run bot without `group-ids` of `[[bot.slack]]` in config.toml
    2. chat `!id` in specified group
    3. check response of `GroupID: XXX`
    4. postscript `group-ids` of `[[bot.slack]]` in config.toml & re-run bot

### setup Bot for Discord

* Setup Slack Bot & Get Bot `TOKEN` : https://discord.com/developers/applications/

* Look up GroupID
    1. run bot without `group-ids` of `[[bot.discord]]` in config.toml
    2. chat `!id` in specified group
    3. check response of `GroupID: XXX`
    4. postscript `group-ids` of `[[bot.discord]]` in config.toml & re-run bot

### Bot needs to support HTTPS separately

This bot run HTTP server, but Webhook configuration required HTTPS in most chat-provider.
Please following the below.

* using HTTPS reverse-proxy server (nginx, Caddy, etc..) & run Bot beside Minecraft server
* using PaaS (Heroku, Google App Engine, etc..)
    * not recommended (RCON connection is not crypted)

