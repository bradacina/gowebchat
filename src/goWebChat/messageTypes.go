package goWebChat

import (
	"encoding/json"
	"log"
)

type MessageType struct {
	Type string
}

// chat message sent from client to server
type ClientChatMessage struct {
	Chat string
}

// chat message sent from server to client
type ServerChatMessage struct {
	Type string
	Chat string
	Name string
}

// server status message sent to client
// contains text that should be displayed on client's screen
type ServerStatusMessage struct {
	Type    string
	Content string
}

// server message containing the list of
// connected clients that are sent to the client
type ServerClientListMessage struct {
	Type    string
	Content string
}

// server message that notifies other clients when a new client
// has connected
type ServerClientJoinMessage struct {
	Type string
	Name string
}

// server message that notifies other clients that a particular
// client has disconnected
type ServerClientPartMessage struct {
	Type string
	Name string
}

func NewServerChatMessage(content string, sender string) ServerChatMessage {
	return ServerChatMessage{Type: "ServerChatMessage", Name: sender, Chat: content}
}

func NewServerStatusMessage(content string) ServerStatusMessage {
	return ServerStatusMessage{Type: "ServerStatusMessage", Content: content}
}

func NewServerClientListMessage(content string) ServerClientListMessage {
	return ServerClientListMessage{Type: "ServerClientListMessage", Content: content}
}

func NewServerClientJoinMessage(name string) ServerClientJoinMessage {
	return ServerClientJoinMessage{Type: "ServerClientJoinMessage", Name: name}
}

func NewServerClientPartMessage(name string) ServerClientPartMessage {
	return ServerClientPartMessage{Type: "ServerClientPartMessage", Name: name}
}

func UnmarshalClientChatMessage(msg []byte) (ClientChatMessage, error) {
	var chatMsg ClientChatMessage
	err := json.Unmarshal(msg, &chatMsg)
	if err != nil {
		log.Println("Error unmarshalling ClientChatMessage", err)
	}
	return chatMsg, err
}
