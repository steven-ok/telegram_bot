package route

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tg_bots/route/bind"
	"tg_bots/service/topic"
)

var botApi *tgbotapi.BotAPI

func NewRoute(api *tgbotapi.BotAPI) *bind.Map {
	botApi = api
	return bind.NewRouteMap(routes)
}

func routes(routeMap *bind.Map) {
	topicService := topic.NewTopicService(botApi)

	routeMap.AddCommandRoute("topic", topicService.CreateTopic)
	routeMap.AddRegularExpression(`^#(\d+)`, topicService.ReplyToTopic)
	routeMap.AddCommandRoute("reply", topicService.ReplyToTopic)
}
