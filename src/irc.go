package src

import (
	//"fmt"
	"net"
	"time"

	"gopkg.in/irc.v3"
)

var userID 		string
var client 		*irc.Client


func join() (string) {
	// dial in
	connect := "localhost:6667"
	conn, err := net.Dial("tcp", connect)
	if(err != nil) {return "Couldn't connect to "+connect+": "+err.Error()}
	config := irc.ClientConfig{
		Nick: userID,
		User: "terminal_wars_"+userID,
		Handler: irc.HandlerFunc(ircHandler),
	}
	client = irc.NewClient(conn, config)
	// Wait two seconds to make absolutely sure the client is initialized
	time.Sleep(2 * time.Second) 
	err = client.Run()
	if(err != nil) {
		return "Disconnected from "+connect+": "+err.Error()
		inIRC = false
	}
	return ""
}

func ircHandler(c *irc.Client, m *irc.Message) {
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