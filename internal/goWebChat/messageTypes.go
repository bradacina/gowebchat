package goWebChat

import (
	"encoding/json"
	"log"
)

// MessageType is the type of message being sent.
type MessageType struct {
	Type string
}

// ClientChatMessage represents a message sent from client to server.
type ClientChatMessage struct {
	Chat string
}

// ClientChangeNameMessage represents a change name message sent from client to server.
type ClientChangeNameMessage struct {
	NewName string
}

// ClientPongMessage represents a pong message sent from client to server.
type ClientPongMessage struct {
	Type    string
	Payload int
}

// ServerPingMessage represents a ping message sent from server to client.
type ServerPingMessage struct {
	Type    string
	Payload int
}

// ServerChatMessage represents a chat message sent from server to client.
type ServerChatMessage struct {
	Type string
	Chat string
	Name string
}

// ServerStatusMessage represents a server status message sent to client.
// It contains text that should be displayed on client's screen.
type ServerStatusMessage struct {
	Type    string
	Content string
}

// ServerClientListMessage represents a server message containing the list of
// connected clients that are sent to the client.
type ServerClientListMessage struct {
	Type    string
	Content string
}

// ServerClientJoinMessage represents a server message that notifies other
// clients when a new client has connected.
type ServerClientJoinMessage struct {
	Type string
	Name string
}

// ServerClientPartMessage represents a server message that notifies other
// clients that a particular client has disconnected.
type ServerClientPartMessage struct {
	Type string
	Name string
}

// ServerSetNameMessage represents a server message that notifies the client
// of their new name.
type ServerSetNameMessage struct {
	Type    string
	NewName string
}

// ServerChangeNameMessage represents a server message that notifies other
// clients that a particular client has changed his/her name.
type ServerChangeNameMessage struct {
	Type    string
	NewName string
	OldName string
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

func NewServerSetNameMessage(name string) ServerSetNameMessage {
	return ServerSetNameMessage{Type: "ServerSetName", NewName: name}
}

func NewServerChangeNameMessage(oldName string, newName string) ServerChangeNameMessage {
	return ServerChangeNameMessage{Type: "ServerChangeName", NewName: newName, OldName: oldName}
}

func NewServerPingMessage(payload int) ServerPingMessage {
	return ServerPingMessage{Type: "ServerPingMessage", Payload: payload}
}

func UnmarshalClientChatMessage(msg []byte) (ClientChatMessage, error) {
	var chatMsg ClientChatMessage
	err := json.Unmarshal(msg, &chatMsg)
	if err != nil {
		log.Println("Error unmarshalling ClientChatMessage", err)
	}
	return chatMsg, err
}

func UnmarshalClientChangeNameMessage(msg []byte) (ClientChangeNameMessage, error) {
	var changeNameMsg ClientChangeNameMessage
	err := json.Unmarshal(msg, &changeNameMsg)
	if err != nil {
		log.Println("Error unmarshalling ClientChatMessage", err)
	}
	return changeNameMsg, err
}

func UnmarshalClientPongMessage(msg []byte) (ClientPongMessage, error) {
	var clientPongMessage ClientPongMessage
	err := json.Unmarshal(msg, &clientPongMessage)
	if err != nil {
		log.Println("Error unmarshalling ClientPongMessage", err)
	}
	return clientPongMessage, err
}
