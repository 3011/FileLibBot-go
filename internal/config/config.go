package config

import (
	"log"

	"github.com/jinzhu/configor"
)

var Config = struct {
	BotToken       string `required:"true"`
	ForwardGroupID int64  `required:"true"`
	DBFileName     string `default:"FileLibBot.db"`
}{}

func Init() {
	err := configor.Load(&Config, "config.yml")
	if err != nil {
		log.Panic(err.Error())
	}
}
