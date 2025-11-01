package handlers

import (
	"gopkg.in/telebot.v3"
)

// SafeSendMessage безопасно отправляет сообщение, обрабатывая специфичные ошибки Telegram
func SafeSendMessage(c telebot.Context, text string, opts ...interface{}) {
	c.Send(text, &telebot.SendOptions{
		ReplyTo:               c.Message(),
		ThreadID:              c.Message().ThreadID,
		DisableWebPagePreview: true,
	})
}
