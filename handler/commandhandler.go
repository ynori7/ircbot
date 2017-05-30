package handler

import (
	"github.com/ynori7/go-irc/client"
	"github.com/ynori7/ircbot/ircconfig"
)

type CommandHandler struct {
	config     ircconfig.IrcConfig
	mutedUsers []string //todo: make this a map and add a mutex
}

func NewCommandHandler(config ircconfig.IrcConfig) CommandHandler {
	return CommandHandler{
		config:     config,
		mutedUsers: make([]string, 0),
	}
}

func (h CommandHandler) UnmuteUser(conn client.Client, nick, location string) {
	conn.SetMode(location, "+v", nick)
	//todo remove from list
}

func (h CommandHandler) MuteUser(conn client.Client, command, location string) {
	//todo
}
