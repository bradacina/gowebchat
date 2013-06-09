package main

import (
	"encoding/json"
	"goWebChat"
	"log"
)

func handleMessage(msg []byte, client goWebChat.Client) {

	var messageType goWebChat.MessageType

	json.Unmarshal(msg, &messageType)
	log.Println("MessageType ", messageType)

	if messageType == (goWebChat.MessageType{}) {
		return
	}

	switch messageType.Type {
	case "ClientChat":
		chatMsg, err := goWebChat.UnmarshalClientChatMessage(msg)

		if err != nil || len(chatMsg.Chat) == 0 {
			return
		}

		var outboundChatMsg = goWebChat.NewServerChatMessage(chatMsg.Chat, client.Name)
		outboundRaw, err := json.Marshal(outboundChatMsg)

		if err != nil {
			log.Println("Error marshaling message to JSON", outboundChatMsg)
			return
		}

		go BroadcastToAll(outboundRaw)
	}
}
