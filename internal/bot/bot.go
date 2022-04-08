package bot

import (
	"log"

	"github.com/3011/FileLibBot-go/internal/callbackquery"
	"github.com/3011/FileLibBot-go/internal/config"
	"github.com/3011/FileLibBot-go/internal/message"
	"github.com/3011/FileLibBot-go/internal/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init(bot *tgbotapi.BotAPI) {
	sqlite.InitDB()
	message.Init(bot)
	callbackquery.Init(bot)

}

func BotStart() {
	botToken := config.Config.BotToken
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	Init(bot)
	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callbackquery.Handle(update.CallbackQuery)
		}

		if update.InlineQuery != nil {
		}

		if update.Message != nil {
			message.Handle(update.Message)
		}
	}
}
