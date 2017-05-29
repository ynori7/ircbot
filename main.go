package main

import (
	"io/ioutil"
	"os"
	"log"
	"errors"

	"github.com/ynori7/ircbot/ircconfig"
	"github.com/ynori7/go-irc/client"
	"github.com/ynori7/ircbot/handler"
)

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
	conn, err := client.NewConnection(config.ConnectionString, config.UseSSL, config.Nick)
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.NewMessageHandler(config)

	conn.Listen(handler.Handle)

}
