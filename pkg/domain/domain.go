package domain

import "github.com/line/line-bot-sdk-go/linebot"

type LineConfig struct {
	ChannelSecret string
	ChannelToken  string
	GroupIDs      []string
}

type LineClient struct {
	Client   *linebot.Client
	GroupIDs []string
}

/* TODO
type SlackClientConfig struct {
	Client   *linebot.Client
	GroupIDs []string
}
type DiscordClientConfig struct {
	Client   *linebot.Client
	GroupIDs []string
}
*/

type Entity struct {
	//ログインしたことのあるすべてのユーザ
	AllUsers []User
	//ログインしてるユーザ
	LoginUsers []User
	//whitelistなユーザ名
	WhitelistUsernames []string
}

type User struct {
	Name     string
	Health   float32
	XpLevel  uint
	Position Position
	Biome    string // Minecraft 1.16~
}

type Position struct {
	X float32
	Y float32
	Z float32
}
