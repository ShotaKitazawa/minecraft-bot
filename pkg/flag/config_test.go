package flag

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bindAddrForTest             = `0.0.0.0`
	bindPortForTest             = 8080
	minecraftHostnameForTest    = `test`
	botLINEEndpointForTest      = `/linebot`
	botLINEChannelSecretForTest = `channel-secret`
	botLINEChannelTokenForTest  = `channel-token`
	botLINEGroupIDForTest       = `group-id`
	rconHostForTest             = `127.0.0.1`
	rconPortForTest             = 25575
	rconPasswordForTest         = `rcon-password`
	sharedmemModeForTest        = `redis`
	sharedmemRedisHostForTest   = `127.0.0.1`
	sharedmemRedisPortForTest   = 6379
)

func TestConfig(t *testing.T) {
	t.Run(`ScanConfig()`, func(t *testing.T) {
		t.Run(`normal`, func(t *testing.T) {
			conf, err := ScanConfig(NewBuffer(bindAddrForTest, bindPortForTest, minecraftHostnameForTest, botLINEEndpointForTest, botLINEChannelSecretForTest, botLINEChannelTokenForTest, botLINEGroupIDForTest, rconHostForTest, rconPortForTest, rconPasswordForTest, sharedmemModeForTest, sharedmemRedisHostForTest, sharedmemRedisPortForTest))

			assert.Nil(t, err)
			assert.Equal(t, bindAddrForTest, conf.BindAddr)
			assert.Equal(t, bindPortForTest, conf.BindPort)
			assert.Equal(t, minecraftHostnameForTest, conf.MinecraftHostname)
			assert.Equal(t, botLINEEndpointForTest, conf.Bot.LINEConfig.Endpoint)
			assert.Equal(t, botLINEChannelSecretForTest, conf.Bot.LINEConfig.ChannelSecret)
			assert.Equal(t, botLINEChannelTokenForTest, conf.Bot.LINEConfig.ChannelToken)
			assert.Equal(t, botLINEGroupIDForTest, conf.Bot.LINEConfig.GroupIDs)
			assert.Equal(t, rconHostForTest, conf.Rcon.Host)
			assert.Equal(t, rconPortForTest, conf.Rcon.Port)
			assert.Equal(t, rconPasswordForTest, conf.Rcon.Password)
			assert.Equal(t, sharedmemModeForTest, conf.SharedMem.Mode)
			assert.Equal(t, sharedmemRedisHostForTest, conf.SharedMem.RedisConfig.Host)
			assert.Equal(t, sharedmemRedisPortForTest, conf.SharedMem.RedisConfig.Port)
		})
		t.Run(`abnormal`, func(t *testing.T) {
			buffer := bytes.NewBufferString(`
!!invalid format!!
`)
			_, err := ScanConfig(buffer)
			assert.NotNil(t, err)
		})
	})
	t.Run(`ValidateConfig()`, func(t *testing.T) {
		t.Run(`valid`, func(t *testing.T) {
			conf := &Config{
				MinecraftHostname: minecraftHostnameForTest,
				Bot: BotConfig{
					LINEConfig: LINEConfig{
						Endpoint:      botLINEEndpointForTest,
						ChannelSecret: botLINEChannelSecretForTest,
						ChannelToken:  botLINEChannelTokenForTest,
					},
				},
				Rcon: RconConfig{
					Password: rconPasswordForTest,
				},
				SharedMem: SharedMemConfig{
					Mode: sharedmemModeForTest,
				},
			}
			err := ValidateConfig(conf)
			assert.Nil(t, err)
		})
		t.Run(`invalid: [minecraft-hostname] is empty`, func(t *testing.T) {
			conf := &Config{
				Bot: BotConfig{
					LINEConfig: LINEConfig{
						Endpoint:      botLINEEndpointForTest,
						ChannelSecret: botLINEChannelSecretForTest,
						ChannelToken:  botLINEChannelTokenForTest,
					},
				},
				Rcon: RconConfig{
					Password: rconPasswordForTest,
				},
				SharedMem: SharedMemConfig{
					Mode: sharedmemModeForTest,
				},
			}
			err := ValidateConfig(conf)
			assert.NotNil(t, err)
		})
		t.Run(`invalid: [log-level] is invalid`, func(t *testing.T) {
			conf := &Config{
				LogLevel:          `invalid`,
				MinecraftHostname: minecraftHostnameForTest,
				Bot: BotConfig{
					LINEConfig: LINEConfig{
						Endpoint:      botLINEEndpointForTest,
						ChannelSecret: botLINEChannelSecretForTest,
						ChannelToken:  botLINEChannelTokenForTest,
					},
				},
				Rcon: RconConfig{
					Password: rconPasswordForTest,
				},
				SharedMem: SharedMemConfig{
					Mode: sharedmemModeForTest,
				},
			}
			err := ValidateConfig(conf)
			assert.NotNil(t, err)
		})
		t.Run(`invalid: [bot.line] is empty`, func(t *testing.T) {
			conf := &Config{
				MinecraftHostname: minecraftHostnameForTest,
				Bot: BotConfig{
					LINEConfig: LINEConfig{},
				},
				Rcon: RconConfig{
					Password: rconPasswordForTest,
				},
				SharedMem: SharedMemConfig{
					Mode: sharedmemModeForTest,
				},
			}
			err := ValidateConfig(conf)
			assert.NotNil(t, err)
		})
		t.Run(`invalid: [bot.rcon.password] is empty`, func(t *testing.T) {
			conf := &Config{
				MinecraftHostname: minecraftHostnameForTest,
				Bot: BotConfig{
					LINEConfig: LINEConfig{
						Endpoint:      botLINEEndpointForTest,
						ChannelSecret: botLINEChannelSecretForTest,
						ChannelToken:  botLINEChannelTokenForTest,
					},
				},
				Rcon: RconConfig{},
				SharedMem: SharedMemConfig{
					Mode: `local`,
				},
			}
			err := ValidateConfig(conf)
			assert.NotNil(t, err)
		})
		t.Run(`invalid: [bot.sharedmem.mode] is invalid`, func(t *testing.T) {
			conf := &Config{
				MinecraftHostname: minecraftHostnameForTest,
				Bot: BotConfig{
					LINEConfig: LINEConfig{
						Endpoint:      botLINEEndpointForTest,
						ChannelSecret: botLINEChannelSecretForTest,
						ChannelToken:  botLINEChannelTokenForTest,
					},
				},
				Rcon: RconConfig{
					Password: rconPasswordForTest,
				},
				SharedMem: SharedMemConfig{
					Mode: `invalid`,
				},
			}
			err := ValidateConfig(conf)
			assert.NotNil(t, err)
		})
	})
}

func NewBuffer(
	bindAddr string,
	bindPort int,
	minecraftHostname string,
	botLINEEndpoint, botLINEChannelSecret, botLINEChannelToken, botLINEGroupID string,
	rconHost string,
	rconPort int,
	rconPassword string,
	sharedmemMode string,
	sharedmemRedisHost string,
	sharedmemRedisPort int,
) *bytes.Buffer {
	return bytes.NewBufferString(fmt.Sprintf(`
bind-addr = "%s"
bind-port = %d
minecraft-hostname = "%s"

[bot.line]
endpoint = "%s"
channel-secret = "%s"
channel-token = "%s"
group-ids = "%s"

[rcon]
host = "%s"
port = %d
password = "%s"

[sharedmem]
mode = "%s"

[sharedmem.redis]
host = "%s"
port = %d
`, bindAddr, bindPort, minecraftHostname, botLINEEndpoint, botLINEChannelSecret, botLINEChannelToken, botLINEGroupID, rconHost, rconPort, rconPassword, sharedmemMode, sharedmemRedisHost, sharedmemRedisPort))
}
