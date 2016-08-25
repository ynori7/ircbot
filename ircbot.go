package main

import (
	"bufio"
	"io/ioutil"
	"fmt"
	"os"
	"log"
	"strings"
	"github.com/ynori7/ircbot/ircconfig"
	"github.com/ynori7/ircbot/ircutil"
	"errors"
)

func HandleMessage(conn ircutil.IrcConnection, message string) {
	line := conn.ParseLine(message)
	sender := ircutil.ParseUserString(line.Sender)

	if line.Type == "PING" {
		conn.Pong(line.Sender)
	}
	//001 appears when we've connected and the server starts talking to us
	if line.Type == "001" || (line.Type == "KICK" && line.Message == conn.Config.Nick){
		conn.JoinChannel(conn.Config.Channels[0]) //TODO: rejoin the channel I was kicked from
	}
	if line.Type == "JOIN" && sender.Nick != conn.Config.Nick {
		conn.SendMessage("hey " + sender.Nick, line.Location)
	}
	if line.Type == "PRIVMSG" && strings.Contains(line.Message, "hello "+conn.Config.Nick) {
		loc := line.Location

		if(line.Location == conn.Config.Nick) {
			loc = sender.Nick
		}

		conn.SendMessage("hi", loc)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("You must specify the path to the config file."))
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var config ircconfig.IrcConfig
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}

	conn := ircutil.IrcConnection{Config: config}
	err = conn.Connect()

	if err != nil {
		log.Fatal(err)
	}

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