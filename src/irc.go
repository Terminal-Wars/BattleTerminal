package src

import (
	//"fmt"
	"net"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/irc.v3"
)

var userID 		string
var tempToken 	string	= ""
var client 		*irc.Client


func join() (string) {
	// dial in
	connect := "play.ioi-xd.net:6667"
	conn, err := net.Dial("tcp", connect)
	if(err != nil) {return "Couldn't connect to "+connect+": "+err.Error()}
	// if our temporary token is blank
	if(tempToken == "") {
		// make one by hashing our IP
		//ip := conn.LocalAddr().(*net.UDPAddr).IP
		ip := []byte("1234567")
		bytes, err := bcrypt.GenerateFromPassword(ip, len(ip))
		if(err != nil) {return "Couldn't create a token: "+err.Error()} 
		tempToken = string(bytes)
	}
	config := irc.ClientConfig{
		Nick: userID,
		//User: "terminal_wars_"+tempToken,
		User: userID,
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