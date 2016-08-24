package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/ynori7/ircbot/ircutil"
)

func HandleMessage(conn ircutil.IrcConnection, message string) {
	line := conn.ParseLine(message)
	sender := ircutil.ParseUserString(line.Sender)

	if line.Type == "PING" {
		conn.Pong(line.Sender)
	}
	//001 appears when we've connected and the server starts talking to us
	if line.Type == "001" || (line.Type == "KICK" && line.Message == conn.Nick){
		conn.JoinChannel(conn.IrcChannel)
	}
	if line.Type == "JOIN" && sender.Nick != conn.Nick {
		conn.SendMessage("hey " + sender.Nick, line.Location)
	}
	if line.Type == "PRIVMSG" && strings.Contains(line.Message, "hello "+conn.Nick) {
		loc := line.Location

		if(line.Location == conn.Nick) {
			loc = sender.Nick
		}

		conn.SendMessage("hi", loc)
	}
}

func main() {
	conn := ircutil.IrcConnection{
		Nick: "blorgleflorp",
		IrcChannel: "#lore",
		ConnectionString: "irc.psych0tik.net:6697",
		UseSSL: true,
		InputChannel: make(chan string),
	}
	err := conn.Connect()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	connbuf := bufio.NewReader(conn.Connection)
	
	for{
		str, err := connbuf.ReadString('\n')
		if len(str)>0 {
			fmt.Println(str)
			go HandleMessage(conn, str)
		}
		if err!= nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}