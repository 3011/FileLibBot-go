package callbackquery

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func Init(botAPI *tgbotapi.BotAPI) {
	bot = botAPI
}

func Handle(callbackQuery *tgbotapi.CallbackQuery) {
	if callbackQuery.Message.Chat.ID != callbackQuery.From.ID {
		msg := tgbotapi.NewCallback(callbackQuery.ID, "这不是你的消息")
		bot.Send(msg)
	}

	dataArr := strings.Fields(callbackQuery.Data)

	switch dataArr[0] {
	case "del":
		deleteFile(callbackQuery.From.ID, dataArr[1])
		msg := tgbotapi.NewCallback(callbackQuery.ID, "Deleted!")
		bot.Send(msg)
	case "page":
		if dataArr[1] == "last" {
			pageInt, _ := strconv.Atoi(dataArr[2])
			commandMyFile(callbackQuery.From.ID, pageInt, callbackQuery.Message.MessageID)

		} else {
			pageInt, _ := strconv.Atoi(dataArr[2])
			commandMyFile(callbackQuery.From.ID, pageInt, callbackQuery.Message.MessageID)
		}
	}
}
