package handler

import (
	"strings"
	"time"

	"github.com/ynori7/go-irc/client"
	"github.com/ynori7/go-irc/model"
	"github.com/ynori7/ircbot/ircconfig"
	"github.com/ynori7/ircbot/library"
)

type Handler struct {
	config ircconfig.IrcConfig
}

func NewMessageHandler(config ircconfig.IrcConfig) Handler {
	return Handler{
		config: config,
	}
}

/**
 * Performs the designated action according to the content of the message received.
 */
func (h Handler) Handle(conn client.Client, message model.Message) {
	if message.Type == client.PING {
		conn.Pong(message.Message)
	}

	if message.Type == "001" { //001 appears when we've connected and the server starts talking to us
		conn.SendMessage("identify "+h.config.Password, "NickServ")

		for _, ch := range h.config.Channels { //join all the channels in the config
			conn.JoinChannel(ch)
		}
	}

	if message.Type == client.KICK && message.Message == conn.Nick {
		conn.JoinChannel(message.Location) //rejoin the channel I was kicked from
	}

	if message.Type == client.JOIN && message.Sender.Nick != conn.Nick { //Greet user who joined channel
		if h.in_array(h.config.ModeratedChannels, message.Location) {
			conn.SetMode(message.Location, "+v", message.Sender.Nick)
		}

		go func() { //to avoid sending the message so fast that the user doesn't notice it
			time.Sleep(500 * time.Millisecond)
			conn.SendMessage(h.config.GetRandomGreeting()+" "+message.Sender.Nick, message.Location)
		}()
	}

	if message.Type == client.PRIVMSG {
		isAdmin := h.in_array(h.config.Admins, message.Sender.Nick)

		if strings.HasPrefix(message.Message, h.config.Nick+":") {
			h.doCommand(conn, message, isAdmin)
		}

		h.doConversation(conn, message, isAdmin)
	}
}

func (h Handler) doCommand(conn client.Client, message model.Message, senderIsAdmin bool) {
	//todo
}

/**
 * Handles conversational type messages like talking to other users.
 */
func (h Handler) doConversation(conn client.Client, message model.Message, senderIsAdmin bool) {
	location := message.Location
	//Handle the case when user is talking to me in private message, not in channel
	if message.Location == conn.Nick {
		location = message.Sender.Nick
	}

	if strings.Contains(message.Message, conn.Nick) { //someone is talking to me or about me
		words := strings.Fields(message.Message)

		//respond to greetings
		if h.in_array(h.config.Greetings, strings.ToLower(words[0])) {
			conn.SendMessage(h.config.GetRandomGreeting(), location)
		}
	}

	if strings.Contains(message.Message, "github.com") {
		githubResponse := library.HandleGithubLink(message.Message)
		if githubResponse != "" {
			conn.SendMessage(githubResponse, location)
		}
	}
}

/**
 * Returns true if needle occurs in haystack, otherwise false.
 * Not sure why there isn't already a function for this.
 */
func (h Handler) in_array(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}
