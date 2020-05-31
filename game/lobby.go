package game

import (
	"duel-masters/server"
	"encoding/json"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	messageBufferSize int = 5
)

var lobby = &Lobby{}

var messages = make([]server.LobbyChatMessage, 0)
var messagesMutex = &sync.Mutex{}

var subscribers = make([]*server.Socket, 0)
var subscribersMutex = &sync.Mutex{}

// Lobby struct is used to create a Hub that can parse messages from the websocket server
type Lobby struct{}

// GetLobby returns a reference to the lobby
func GetLobby() *Lobby {
	return lobby
}

// Broadcast sends a message to all subscribed sockets
func Broadcast(msg interface{}) {
	subscribersMutex.Lock()
	defer subscribersMutex.Unlock()

	for _, subscriber := range subscribers {
		go subscriber.Send(msg)
	}
}

// Parse websocket messages
func (l *Lobby) Parse(s *server.Socket, data []byte) {

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from parsing a message in lobby. %v", r)
		}
	}()

	var message server.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return
	}

	switch message.Header {

	case "subscribe":
		{
			subscribersMutex.Lock()
			defer subscribersMutex.Unlock()

			for _, subscriber := range subscribers {
				if subscriber == s {
					return
				}
			}

			subscribers = append(subscribers, s)

			// Send chat messages
			s.Send(server.LobbyChatMessages{
				Header:   "chat",
				Messages: messages,
			})

		}

	case "chat":
		{

			var msg struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			messagesMutex.Lock()
			defer messagesMutex.Unlock()

			if len(messages) >= messageBufferSize {
				_, messages = messages[0], messages[1:]
			}

			chatMsg := server.LobbyChatMessage{
				Username:  s.User.Username,
				Color:     "orange",
				Message:   msg.Message,
				Timestamp: int(time.Now().Unix()),
			}

			toBroadcast := server.LobbyChatMessages{
				Header:   "chat",
				Messages: []server.LobbyChatMessage{chatMsg},
			}

			messages = append(messages, chatMsg)

			Broadcast(toBroadcast)

		}

	}

}
