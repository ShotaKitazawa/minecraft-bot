package i18n

import (
	"golang.org/x/text/message"
)

var T *message.Printer

func init() {
	// default
	initEnglish()

	// when LANG=ja.*
	initJapanese()
}

var (
	MessageError             = `MessageError`
	MessageHelp              = `MessageHelp`
	MessageInvalidArguments  = `MessageInvalidArguments`
	MessageMemberJoined      = `MessageMemberJoined: %s`
	MessageNoLoginUserExists = `MessageNoLoginUserExists`
	MessageNoSuchCommand     = `MessageNoSuchCommand`
	MessageNoUserExists      = `MessageNoUserExists`
	MessageSentMessage       = `MessageSentMessage: %s`
	MessageUserIncorrect     = `MessageUserIncorrect: %s`
	MessageUsersLogin        = `MessageUsersLogin: %v`
	MessageUsersLogout       = `MessageUsersLogout: %v`
	MessageWhitelistAdd      = `MessageWhitelistAdd: %s`
	MessageWhitelistRemove   = `MessageWhitelistRemove: %s`
)
