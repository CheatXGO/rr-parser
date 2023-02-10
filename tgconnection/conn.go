package tgconnection

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SetWebhook(bot *tgbotapi.BotAPI, token string, ip string, port string, cert string) error {
	webHookInfo := tgbotapi.NewWebhookWithCert(fmt.Sprintf("https://%s:%s/%s", ip, port, token), cert)
	_, err := bot.SetWebhook(webHookInfo)
	if err != nil {
		return err
	}
	return nil
}
