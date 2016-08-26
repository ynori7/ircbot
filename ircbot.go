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
	if line.Type == "KICK" {
		fmt.Printf("%v", line.Message)
	}
	if line.Type == "KICK" && line.Message == conn.Config.Nick {
		conn.JoinChannel(line.Location) //rejoin the channel you were kicked from
	}

	if line.Type == "JOIN" && line.Sender.Nick != conn.Config.Nick {
		conn.SendMessage("hey " + line.Sender.Nick, line.Location)
	}
	if line.Type == "PRIVMSG" {
		Conversation(conn, line)
	}
}

func Conversation(conn ircutil.IrcConnection, line ircutil.IrcMessage) {
	location := line.Location
	if(line.Location == conn.Config.Nick) {
		location = line.Sender.Nick
	}

	if strings.Contains(line.Message, "hello "+conn.Config.Nick) {
		conn.SendMessage("hi", location)
	}
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
			go HandleMessage(conn, str)
		}
		if err!= nil {
			log.Fatal(err)
		}
	}

}