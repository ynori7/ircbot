package handler

import (
	"log"
	"strings"
	"time"

	"github.com/ynori7/go-irc/client"
	"github.com/ynori7/go-irc/model"
	"github.com/ynori7/ircbot/ircconfig"
	"github.com/ynori7/ircbot/library"
	"github.com/ynori7/ircbot/service"
)

type MessageHandler struct {
	config       ircconfig.IrcConfig
	voiceService service.VoiceService
}

func NewMessageHandler(config ircconfig.IrcConfig, voiceService service.VoiceService) MessageHandler {
	return MessageHandler{
		config:       config,
		voiceService: voiceService,
	}
}

/**
 * Performs the designated action according to the content of the message received.
 */
func (h MessageHandler) Handle(conn *client.Client, message model.Message) {
        log.Println(message.Raw)

	if message.Type == client.PING {
		conn.Pong(message.Message)
	}

	if message.Type == "001" { //001 appears when we've connected and the server starts talking to us
		if h.config.Password != "" {
			conn.SendMessage("identify "+h.config.Password, "NickServ")
		}

		for _, ch := range h.config.Channels { //join all the channels in the config
			conn.JoinChannel(ch)
		}
	}

	if message.Type == client.KICK && message.Message == conn.Nick {
		conn.JoinChannel(message.Location) //rejoin the channel I was kicked from
	}

	if message.Type == client.JOIN && message.Sender.Nick != conn.Nick { //Greet user who joined channel
		if h.in_array(h.config.ModeratedChannels, message.Location) &&
			!h.voiceService.IsMuted(message.Sender.Nick, message.Location) {

			h.voiceService.GiveVoice(message.Sender.Nick, message.Location)
		}

		go func() { //to avoid sending the message so fast that the user doesn't notice it
			time.Sleep(500 * time.Millisecond)
			conn.SendMessage(h.config.GetRandomGreeting()+" "+message.Sender.Nick, message.Location)
		}()
	}

	if message.Type == client.PRIVMSG {
		isAdmin := h.in_array(h.config.Admins, message.Sender.Nick)
		isCommand := false

		if strings.HasPrefix(message.Message, h.config.Nick+":") {
			isCommand = h.doCommand(conn, message, isAdmin)
		}

		if !isCommand {
			h.doConversation(conn, message, isAdmin)
		}
	}
}

/**
 * Handles explicit commands issued to the bot. Commands are prefixed with the bot's name and a colon
 * returns true if there was really a command in the message
 */
func (h MessageHandler) doCommand(conn *client.Client, message model.Message, senderIsAdmin bool) bool {
	commandString := strings.Trim(strings.TrimPrefix(message.Message, h.config.Nick+":"), " ")

	if senderIsAdmin && strings.HasPrefix(commandString, service.MUTE_PREFIX) {
		nick := strings.Trim(strings.TrimPrefix(commandString, service.MUTE_PREFIX), " ")
		if !strings.Contains(nick, " ") {
			h.voiceService.MuteUser(nick, message.Location)
			return true
		}
	}
	if senderIsAdmin && strings.HasPrefix(commandString, service.UNMUTE_PREFIX) {
		nick := strings.Trim(strings.TrimPrefix(commandString, service.UNMUTE_PREFIX), " ")
		if !strings.Contains(nick, " ") {
			h.voiceService.UnmuteUser(nick, message.Location)
			return true
		}
	}

	return false
}

/**
 * Handles conversational type messages like talking to other users.
 */
func (h MessageHandler) doConversation(conn *client.Client, message model.Message, senderIsAdmin bool) {
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
func (h MessageHandler) in_array(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}
