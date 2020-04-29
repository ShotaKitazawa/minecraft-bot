package i18n

import (
	"os"
	"regexp"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func initJapanese() {

	if matched, _ := regexp.Match(`ja_JP.*`, []byte(os.Getenv("LANG"))); !matched {
		return
	}

	T = message.NewPrinter(language.Japanese)

	message.SetString(language.Japanese, MessageError, `内部エラーが発生しました`)

	message.SetString(language.Japanese, MessageHelp, `
/help
ヘルプメッセージを表示します

/list
ログイン中のユーザ一覧を表示します

/title hoge
Minecraftのゲーム画面に hoge と表示されます

/whitelist list
ホワイトリストを表示します

/whitelist add hoge
ユーザ hoge をホワイトリストに追加します

/whitelist delete hoge
ユーザ hoge をホワイトリストから削除します
`)

	message.SetString(language.Japanese, MessageInvalidArguments, `引数が間違っています`)

	message.SetString(language.Japanese, MessageMemberJoined, `
ようこそ！ 以下の手順でセットアップをしてください。
1. このトークルームで "/whitelist add ${MINECRAFT_ACCOUNT_ID}" と発言
2. Minecraft を起動して、"%s" サーバに参加
`)

	message.SetString(language.Japanese, MessageNoLoginUserExists, `ログイン中のユーザは存在しません`)

	message.SetString(language.Japanese, MessageNoSuchCommand, `コマンドが存在しません`)

	message.SetString(language.Japanese, MessageNoUserExists, `ユーザが存在しません`)

	message.SetString(language.Japanese, MessageSentMessage, `%s に送信しました`)

	message.SetString(language.Japanese, MessageUserIncorrect, `ユーザ指定が間違っています: %s`)

	message.SetString(language.Japanese, MessageUsersLogin, `ユーザがログインしました: %v`)

	message.SetString(language.Japanese, MessageUsersLogout, `ユーザがログアウトしました: %v`)

	message.SetString(language.Japanese, MessageWhitelistAdd, `ユーザをホワイトリストに追加しました: %s`)

	message.SetString(language.Japanese, MessageWhitelistRemove, `ユーザをホワイトリストから削除しました: %s`)
}
