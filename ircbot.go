package main

import (
	"bufio"
	"io/ioutil"
	"fmt"
	"os"
	"log"
	"strings"
	"errors"

	"github.com/ynori7/ircbot/ircconfig"
	"github.com/ynori7/ircbot/ircutil"
)

/**
 * Performs the designated action according to the content of the message received.
 */
func HandleMessage(conn ircutil.IrcConnection, message string) {
	line := conn.ParseLine(message)

	if line.Type == "PING" {
		conn.Pong(line.Message)
	}

	if line.Type == "001" { //001 appears when we've connected and the server starts talking to us
		for _, ch := range conn.Config.Channels { //join all the channels in the config
			conn.JoinChannel(ch)
		}
	}

	if line.Type == "KICK" && line.Message == conn.Config.Nick {
		conn.JoinChannel(line.Location) //rejoin the channel I was kicked from
	}

	if line.Type == "JOIN" && line.Sender.Nick != conn.Config.Nick { //Greet user who joined channel
		conn.SendMessage(conn.Config.GetRandomGreeting() + " " + line.Sender.Nick, line.Location)
	}
	if line.Type == "PRIVMSG" {
		Conversation(conn, line)
	}
}

/**
 * Handles conversational type messages like talking to other users.
 */
func Conversation(conn ircutil.IrcConnection, line ircutil.IrcMessage) {
	location := line.Location
	//Handle the case when user is talking to me in private message, not in channel
	if(line.Location == conn.Config.Nick) {
		location = line.Sender.Nick
	}

	if strings.Contains(line.Message, conn.Config.Nick) { //someone is talking to me or about me
		words := strings.Fields(line.Message)

		//respond to greetings
		if in_array(conn.Config.Greetings, strings.ToLower(words[0])){
			conn.SendMessage(conn.Config.GetRandomGreeting(), location)
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
	var config ircconfig.IrcConfig
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}

	//Connect
	conn := ircutil.IrcConnection{Config: config}
	err = conn.Connect()

	if err != nil {
		log.Fatal(err)
	}

	//Start reading from the connection
	connbuf := bufio.NewReader(conn.Connection)
	for{
		str, err := connbuf.ReadString('\n')
		if len(str)>0 {
			fmt.Println(str)
			go HandleMessage(conn, str) //handle message asynchronously so we can go back to listening
		}
		if err!= nil {
			log.Fatal(err)
		}
	}

}