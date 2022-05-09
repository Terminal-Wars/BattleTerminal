package src

import (
	"fmt"
	"net"
	"errors"

	"gopkg.in/irc.v3"
)

var userID 		string
var client 		*irc.Client


func join() (error) {
	// dial in
	connect := "play.ioi-xd.net:6667"
	conn, err := net.Dial("tcp", connect)
	if(err != nil) {
		return errors.New("Couldn't connect to "+connect+": "+err.Error())
	}
	config := irc.ClientConfig{
		Nick: userID,
		User: "terminal_wars_"+userID,
		Handler: irc.HandlerFunc(ircHandler),
	}
	client = irc.NewClient(conn, config)
	// Wait to make absolutely sure the client is initialized
	for(client == nil) {}
	go func() {
		err = client.Run()
		sendToTextarea("Disconnected from "+connect+": "+err.Error())
		inIRC = false
	}()
	return nil
}

func ircHandler(c *irc.Client, m *irc.Message) {
	fmt.Println(m.String())
	sendToTextarea(m.String())
	if m.Command == "001" {
		// 001 is a welcome event, so we join channels there
		c.Write("JOIN #testing")
	} else if m.Command == "PRIVMSG" && c.FromChannel(m) {
		// Create a handler on all messages.
		c.WriteMessage(&irc.Message{
			Command: "PRIVMSG",
			Params: []string{
				m.Params[0],
				m.Trailing(),
			},
		})
	}
}