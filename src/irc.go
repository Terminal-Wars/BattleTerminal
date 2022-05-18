package src

import (
	"net"
	"errors"

	"gopkg.in/irc.v3"
)

var userID 		string
var client 		*irc.Client


func join() (error) {
	// dial in
	connect := "localhost:6667"
	conn, err := net.Dial("tcp", connect)
	if(err != nil) {
		return errors.New("Couldn't connect to "+connect+": "+err.Error())
	}
	config := irc.ClientConfig{
		Nick: userID,
		User: userID,
		Handler: irc.HandlerFunc(ircHandler),
	}
	client = irc.NewClient(conn, config)
	// Wait to make sure we don't try and do anything when the client is not initialized
	for(client == nil) {}
	go func() {
		err = client.Run()
		sendToTextarea("Disconnected from "+connect+": "+err.Error())
		inIRC = false
	}()
	return nil
}

func joinTestChannel() {
	if(client != nil) {
		client.Write("JOIN #testing")
	}
}

func leave() {
	if(client != nil) {
		client.Write("QUIT")
	}
}

func ircHandler(c *irc.Client, m *irc.Message) {
	sendToTextarea(m.String())
}