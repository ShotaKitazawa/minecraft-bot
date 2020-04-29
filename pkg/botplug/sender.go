package botplug

type BotPluginSender interface {
	SendTextMessage(string) error
}
