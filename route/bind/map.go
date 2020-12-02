package bind

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
	"strings"
	"tg_bots/service/handle"
)

type RoutesBindFun func(routeMap *Map)

type Map struct {
	commandMap        map[string]*Route
	messageMap        map[string]*Route
	regularExpression []*Route
}

func NewRouteMap(routes RoutesBindFun) *Map {
	routeMap := &Map{
		commandMap:        make(map[string]*Route),
		messageMap:        make(map[string]*Route),
		regularExpression: make([]*Route, 0),
	}

	routes(routeMap)

	return routeMap
}

func (routeMap *Map) AddMessageRoute(path string, handler handle.CommandHandler) {
	if path != "" {
		routeMap.messageMap[path] = NewRoute(path, handler, false, false)
	}
}

func (routeMap *Map) AddCommandRoute(path string, handler handle.CommandHandler) {
	if path != "" {
		routeMap.commandMap[path] = NewRoute(path, handler, true, false)
	}
}

func (routeMap *Map) AddRegularExpression(path string, handler handle.CommandHandler) {
	if path != "" {
		routeMap.regularExpression = append(routeMap.regularExpression, NewRoute(path, handler, false, true))
	}
}

// 执行路由动作
func (routeMap *Map) Exec(message *tgbotapi.Message) {
	text := strings.Split(message.Text, " ")

	// 优先执行命令路由
	if command := message.Command(); command != "" {
		if route, ok := routeMap.commandMap[command]; ok {
			//执行对应动作
			route.handle(command, message, text...)
			return
		}
	}

	// 执行消息路由
	if route, ok := routeMap.messageMap[text[0]]; ok {
		//执行对应动作
		route.handle(text[0], message, text...)
		return
	}

	// 遍历正则
	for _, route := range routeMap.regularExpression {
		reg := regexp.MustCompile(route.Command)
		res := reg.FindStringSubmatch(message.Text)
		if len(res) > 0 {
			route.handle(text[0], message, res...)
		}
	}
}
