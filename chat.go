package main

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type MessageOutboundEvent struct {
	Event   string
	Channel string
	Sender  string
	Content string
}

type MessageInboundEvent struct {
	Event   string
	Channel string
	Content string
}

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
	MaxUsers    int
	Mutex       sync.Mutex
}

func MakeChatroom() *Chatroom {
	c := Chatroom{
		MaxChannels: 10,
		MaxUsers:    10,
		Channels:    make(map[string]Channel),
		Users:       make(map[string]User),
	}

	return &c
}

// Check if a sessionID is registered with the chatroom service
func (c *Chatroom) IsRegistered(sessionID string) bool {
	c.Mutex.Lock()
	_, ok := c.Users[sessionID]
	c.Mutex.Unlock()
	return ok
}

// Pre-register the user before the websocket handshake
func (c *Chatroom) Register(sessionID string, displayName string) error {
	c.Mutex.Lock()
	// Check if the sessionId is already registered
	_, ok := c.Users[sessionID]
	if ok {
		c.Mutex.Unlock()
		return errors.New("Session already registered")
	}

	// Check if chatroom is full
	if len(c.Users) == c.MaxUsers {
		c.Mutex.Unlock()
		return errors.New("The chatroom is full")
	}

	// Check if display name has been registered

	c.Users[sessionID] = User{DisplayName: displayName, Connection: nil}

	c.Mutex.Unlock()
	return nil
}

// Accept the user's connection if pre-registered
func (c *Chatroom) Connect(sessionID string, conn *websocket.Conn) error {
	c.Mutex.Lock()
	// Check if the sessionID is already registered
	user, ok := c.Users[sessionID]
	if !ok {
		c.Mutex.Unlock()
		return errors.New("Session not registered")
	}

	// Check if sessionID already has a connection
	if user.Connection != nil {
		c.Mutex.Unlock()
		return errors.New("Session already connected")
	}

	// Check if the chatroom is full
	if len(c.Users) == c.MaxUsers {
		c.Mutex.Unlock()
		return errors.New("The chatroom is full")
	}

	user.Connection = conn
	c.Users[sessionID] = user

	c.Mutex.Unlock()
	return nil
}

// Receive a MessageInboundEvent
func (c *Chatroom) ReceiveMessage(sessionID string, event MessageInboundEvent) error {
	c.Mutex.Lock()
	// Check if the sessionID is already registered
	user, ok := c.Users[sessionID]
	if !ok {
		c.Mutex.Unlock()
		return errors.New("Session not registered")
	}

	// Check if sessionID already has a connection
	if user.Connection == nil {
		c.Mutex.Unlock()
		return errors.New("Session not connected")
	}

	// Check if the channel exists
	ch, ok := c.Channels[event.Channel]
	if !ok {
		c.Mutex.Unlock()
		return errors.New("Channel does not exist")
	}

	// Add the message to the channel's messages
	ch.Messages = append(ch.Messages, Message{Sender: user.DisplayName, Message: event.Content})

	// Update the channel
	c.Channels[event.Channel] = ch

	// Broadcast the message to all connected users
	x := &MessageOutboundEvent{
		Event:   "message",
		Sender:  user.DisplayName,
		Channel: event.Channel,
		Content: event.Content,
	}
	y, err := json.Marshal(x)
	if err != nil {
		c.Mutex.Unlock()
		return err
	}

	for _, user := range c.Users {
		user.Connection.WriteMessage(1, []byte(y))
	}

	c.Mutex.Unlock()
	return nil
}

// Add a channel by name
func (c *Chatroom) AddChannel(name string) error {
	c.Mutex.Lock()
	// Check if we have room for more channels
	if c.MaxChannels == len(c.Channels) {
		c.Mutex.Unlock()
		return errors.New("Maximum number of channels reached")
	}

	// Check if channel name exists
	if _, ok := c.Channels[name]; ok {
		c.Mutex.Unlock()
		return errors.New("A channel with that name already exists")
	}

	c.Channels[name] = Channel{Name: name, Messages: make([]Message, 10)}

	c.Mutex.Unlock()
	return nil
}
