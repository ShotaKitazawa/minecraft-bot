package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func initEnglish() {

	T = message.NewPrinter(language.English)

	message.SetString(language.English, MessageError, `Internal Error`)

	message.SetString(language.English, MessageHelp, `
/help
display help messages

/list
display log in users name

/title hoge
display "hoge" to Minecraft

/whitelist list
display whitelist

/whitelist add hoge
add user "hoge" to whitelist

/whitelist delete hoge
delete user "hoge" from whitelist
`)

	message.SetString(language.English, MessageInvalidArguments, `Invalid arguments`)

	message.SetString(language.English, MessageMemberJoined, `
Welcome! Please set up this process
1. chat "/whitelist add ${MINECRAFT_ACCOUNT_ID}" in this talk-room
2. launch Minecraft & enter server to "%s"
`)

	message.SetString(language.English, MessageNoLoginUserExists, `No login user exists`)

	message.SetString(language.English, MessageNoSuchCommand, `No such command`)

	message.SetString(language.English, MessageNoUserExists, `No user exists`)

	message.SetString(language.English, MessageSentMessage, `sent to %s`)

	message.SetString(language.English, MessageUserIncorrect, `incorrect user specification: %s`)

	message.SetString(language.English, MessageUsersLogin, `%v logged in`)

	message.SetString(language.English, MessageUsersLogout, `%v logged out`)

	message.SetString(language.English, MessageWhitelistAdd, `add %s to whitelist`)

	message.SetString(language.English, MessageWhitelistRemove, `remove %s to whitelist`)

}
