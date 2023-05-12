package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
)

const TaskHashtag = `\#задача`

type Update struct {
	tgbotapi.Update
}

func (update *Update) IsTask() bool {
	result, _ := regexp.MatchString(TaskHashtag, update.Message.Text)
	return result
}

