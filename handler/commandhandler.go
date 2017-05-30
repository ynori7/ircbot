package handler

import (
	"github.com/ynori7/go-irc/client"
	"github.com/ynori7/ircbot/ircconfig"
)

type CommandHandler struct {
	config ircconfig.IrcConfig
}

func NewCommandHandler(config ircconfig.IrcConfig) CommandHandler {
	return CommandHandler{
		config: config,
	}
}

func (h CommandHandler) MuteUser(conn client.Client, command string) {
	//todo
}
