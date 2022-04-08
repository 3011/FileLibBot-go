package message

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func Init(botapi *tgbotapi.BotAPI) {
	bot = botapi
}

func Handle(message *tgbotapi.Message) {
	if commandArguments := message.CommandArguments(); commandArguments != "" {
		if strings.HasPrefix(commandArguments, "startfile") {
			fileUniqueId := strings.Replace(commandArguments, "startfile", "", 1)
			sendFile(message.From.ID, fileUniqueId)
		}
	}

	if message.Text == "/file" {
		commandMyFile(message.From.ID, 0, 0)
	}

	if message.Document != nil {
		newDocumentFile(message)
	}

	if message.Photo != nil {
		newPhotoFile(message)
	}

	if message.Video != nil {
		newVideoFile(message)
	}
}
