package ircutil

import (
	"net"
	"crypto/tls"
	"fmt"
	"strings"
	"regexp"
	"github.com/ynori7/ircbot/ircconfig"
)

type IrcConnection struct {
	Config ircconfig.IrcConfig
	Connection net.Conn
}

type IrcMessage struct {
	Type string
	Sender string
	Location string
	Message string
	Raw string
}

type IrcUser struct {
	Nick string
	Username string
	Host string
	Raw string
}

func (c *IrcConnection) Connect() (err error) {
	if c.Config.UseSSL {
		c.Connection, err = tls.Dial("tcp", c.Config.ConnectionString, &tls.Config{InsecureSkipVerify : true})
	} else {
		c.Connection, err = net.Dial("tcp", c.Config.ConnectionString)
	}

	if err == nil {
		fmt.Fprintf(c.Connection, "USER %s %s %s :%s\r\n", c.Config.Nick, c.Config.Nick, c.Config.Nick, c.Config.Nick)
		fmt.Fprintf(c.Connection, "NICK %s\r\n", c.Config.Nick)
	}

	return err
}

func (c *IrcConnection) SendMessage(msg string, to string) {
	fmt.Fprintf(c.Connection, "PRIVMSG %s :%s\r\n", to, msg)
}

func (c *IrcConnection) JoinChannel(channel string) {
	fmt.Fprintf(c.Connection, "JOIN %s\r\n", channel)
}

func (c *IrcConnection) Pong(server string) {
	fmt.Fprintf(c.Connection, "PONG %s\r\n", server)
}

/*
Samples Messages:
:ynori7!~ynori7@unaffiliated/ynori7 KICK #ynori7 blorgleflorps :blorgleflorps
:blorgleflorps!~blorglefl@2001:4c50:29e:2c00:9084:4b28:8dbd:791 JOIN #ynori7
:wolfe.freenode.net 353 blorgleflorps @ #ynori7 :blorgleflorps @ynori7
:wolfe.freenode.net 366 blorgleflorps #ynori7 :End of /NAMES list.
:ynori7!~ynori7@unaffiliated/ynori7 PRIVMSG #ynori7 :hello blorgleflorps
 */
func (c *IrcConnection) ParseLine(msg string) (IrcMessage) {
	ircMsg := IrcMessage{Raw: msg}

	if strings.HasPrefix(msg, "PING") {
		ircMsg.Type = "PING"
		ircMsg.Sender = strings.Fields(msg)[1]
	} else {
		if strings.HasPrefix(msg, ":") {
			msg = msg[1:]
		}

		tmp := strings.Fields(msg)
		ircMsg.Sender = tmp[0]
		ircMsg.Type = tmp[1]

		//For JOIN there's a : in front
		if strings.HasPrefix(tmp[2], ":") {
			tmp[2] = tmp[2][1:]
		}
		ircMsg.Location = tmp[2]

		if len(tmp) >= 3  && strings.Contains(msg, ":") {
			ircMsg.Message = strings.TrimSpace(strings.SplitAfterN(msg, ":", 2)[1])
		}
	}

	return ircMsg
}

func ParseUserString(userString string) (IrcUser) {
	ircUser := IrcUser{Raw: userString}

	re, err := regexp.Compile(`(.*)!(.*)@(.*)`)

	if err == nil {
		res := re.FindStringSubmatch(userString)

		if len(res) == 4 {
			ircUser.Nick = res[1]
			ircUser.Username = res[2]
			ircUser.Host = res[3]
		}
	}

	return ircUser
}