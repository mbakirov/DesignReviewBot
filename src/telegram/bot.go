package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
)

const Token = "2132831960:asdasdasd"

type Telegram struct {
	Key string
	client *tgbotapi.BotAPI
}

func NewBot(key string, debug bool) *Telegram {
	client, err := tgbotapi.NewBotAPI(key)

	if err != nil {
		log.Fatal(err)
	}

	if debug == true {
		client.Debug = true
	}

	log.Printf("Authorized on account %s", client.Self.UserName)

	return &Telegram{
		client: client,
	}
}

func (t *Telegram) SetWebhook(url string) {
	_, err := t.client.SetWebhook(tgbotapi.NewWebhook(url))

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Webhook url updated: %s", url)
	}
}

func (t *Telegram) GetWebhook() tgbotapi.WebhookInfo {
	info, err := t.client.GetWebhookInfo()

	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return info
}

func (t *Telegram) PullUpdates(updateUrl string) tgbotapi.UpdatesChannel {
	updates := t.client.ListenForWebhook(updateUrl)
	go http.ListenAndServe("0.0.0.0:8001", nil)

	return updates
}