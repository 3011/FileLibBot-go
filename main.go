package main

import (
	"github.com/3011/FileLibBot-go/internal/bot"
	"github.com/3011/FileLibBot-go/internal/config"
)

func main() {
	config.Init()
	bot.BotStart()
}
