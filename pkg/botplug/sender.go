package botplug

type BotSender interface {
	SendTextMessageToChannels(string) error
}
