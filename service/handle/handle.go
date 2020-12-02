package handle

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type CommandHandler func(Command string, message *tgbotapi.Message, parameters ...string)
