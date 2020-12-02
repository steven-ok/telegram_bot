package bind

import (
	_ "tg_bots/service"
	"tg_bots/service/handle"
)

type Route struct {
	Command   string
	handle    handle.CommandHandler
	IsRegular bool
	IsCommand bool
}

func NewRoute(command string, handleFun handle.CommandHandler, isCommand bool, isRegular bool) *Route {
	return &Route{
		Command:   command,
		handle:    handleFun,
		IsRegular: isRegular,
		IsCommand: isCommand,
	}
}
