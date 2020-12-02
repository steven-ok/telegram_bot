package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"tg_bots/common"
	"tg_bots/model"
	"tg_bots/route"
)

func init() {
	common.Init()
	model.Init()
}

func main() {
	//TODO 填写自己的 bot secret
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	route := route.NewRoute(bot)

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		} else {
			log.Printf("(USERID:%d|CHAT_ID:%d)[%s] %s", update.Message.From.ID, update.Message.Chat.ID, update.Message.From.String(), update.Message.Text)
			route.Exec(update.Message)
		}
	}
}
