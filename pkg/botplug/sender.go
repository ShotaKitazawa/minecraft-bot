package botplug

type BotSender interface {
	SendTextMessage(string) error
}
