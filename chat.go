package main

import (
	"errors"

	"github.com/gorilla/websocket"
)

type User struct {
	DisplayName string
	Connection  *websocket.Conn
}

type Message struct {
	Sender  string
	Message string
}

type Channel struct {
	Name     string
	Messages []Message
}

type Chatroom struct {
	Users       map[string]User
	Channels    map[string]Channel
	MaxChannels int
}

func MakeChatroom() Chatroom {
	c := Chatroom{
		MaxChannels: 10,
		Channels:    make(map[string]Channel),
		Users:       make(map[string]User),
	}
	c.Channels["main"] = Channel{Name: "main", Messages: make([]Message, 10)}

	return c
}

// Pre-register the user before the websocket handshake
func (c *Chatroom) Register(sessionID string, displayName string) error {
	// Check if the sessionId is already registered
	_, ok := c.Users[sessionID]
	if ok {
		return errors.New("Session already registered")
	}

	// Check if display name has been registered

	c.Users[sessionID] = User{DisplayName: displayName, Connection: nil}

	return nil
}

// Accept the user's connection if pre-registered
func (c *Chatroom) Connect(sessionID string, conn *websocket.Conn) error {
	// Check if the sessionID is already registered
	user, ok := c.Users[sessionID]
	if !ok {
		return errors.New("Session not registered")
	}

	// Check if sessionID already has a connection
	if user.Connection != nil {
		return errors.New("Session already connected")
	}

	user.Connection = conn
	c.Users[sessionID] = user

	return nil
}

func (c *Chatroom) ReceiveMessage(sessionID string, message string, channel string) error {
	// Check if the sessionID is already registered
	user, ok := c.Users[sessionID]
	if !ok {
		return errors.New("Session not registered")
	}

	// Check if sessionID already has a connection
	if user.Connection == nil {
		return errors.New("Session not connected")
	}

	// Check if the channel exists
	ch, ok := c.Channels[channel]
	if !ok {
		return errors.New("Channel does not exist")
	}

	// Add the message to the channel's messages
	ch.Messages = append(ch.Messages, Message{Sender: user.DisplayName, Message: message})

	// Update the channel
	c.Channels[channel] = ch

	// Broadcast the message to all connected users
	for _, user := range c.Users {
		user.Connection.WriteMessage(1, []byte(message))
	}

	return nil
}
