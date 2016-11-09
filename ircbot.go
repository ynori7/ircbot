package main

import (
	"io/ioutil"
	"os"
	"log"
	"strings"
	"errors"

	"github.com/ynori7/ircbot/ircconfig"
	"github.com/ynori7/go-irc/client"
	"github.com/ynori7/go-irc/model"
	"github.com/ynori7/ircbot/library"
)

var config ircconfig.IrcConfig

/**
 * Performs the designated action according to the content of the message received.
 */
func HandleMessage(conn client.Client, message model.Message) {
	if message.Type == "PING" {
		conn.Pong(message.Message)
	}

	if message.Type == "001" { //001 appears when we've connected and the server starts talking to us
		for _, ch := range config.Channels { //join all the channels in the config
			conn.JoinChannel(ch)
		}
	}

	if message.Type == "KICK" && message.Message == conn.Nick {
		conn.JoinChannel(message.Location) //rejoin the channel I was kicked from
	}

	if message.Type == "JOIN" && message.Sender.Nick != conn.Nick { //Greet user who joined channel
		conn.SendMessage(config.GetRandomGreeting() + " " + message.Sender.Nick, message.Location)
	}
	if message.Type == "PRIVMSG" {
		Conversation(conn, message)
	}
}

/**
 * Handles conversational type messages like talking to other users.
 */
func Conversation(conn client.Client, message model.Message) {
	location := message.Location
	//Handle the case when user is talking to me in private message, not in channel
	if(message.Location == conn.Nick) {
		location = message.Sender.Nick
	}

	if strings.Contains(message.Message, conn.Nick) { //someone is talking to me or about me
		words := strings.Fields(message.Message)

		//respond to greetings
		if in_array(config.Greetings, strings.ToLower(words[0])){
			conn.SendMessage(config.GetRandomGreeting(), location)
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
func in_array(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("You must specify the path to the config file."))
	}

	//Get the config
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}

	//Connect
	conn, err := client.NewConnection(config.ConnectionString, config.UseSSL, config.Nick)
	if err != nil {
		log.Fatal(err)
	}

	conn.Listen(HandleMessage)

}