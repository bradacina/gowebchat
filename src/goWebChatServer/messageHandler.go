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
			log.Println("Could not unmarshall incoming message from client: ", msg)
			return
		}

		var outboundChatMsg = goWebChat.NewServerChatMessage(chatMsg.Chat, client.Name)
		outboundRaw, err := json.Marshal(outboundChatMsg)

		if err != nil {
			log.Println("Error marshaling message to JSON", outboundChatMsg)
			return
		}

		go BroadcastToAll(outboundRaw)
	case "ChangeName":
		changeNameMsg, err := goWebChat.UnmarshalClientChangeNameMessage(msg)

		if err != nil || len(changeNameMsg.NewName) == 0 {
			log.Println("Could not unmarshall incoming message from client: ", msg)
			return
		}

		if client.Name == changeNameMsg.NewName {
			return
		}

		// set new name
		oldName := string([]byte(client.Name))
		newName := ChangeName(client.Name, changeNameMsg.NewName)

		log.Println(oldName, newName, changeNameMsg.NewName)

		SendName(newName, &client)

		// broadcast the name change to everyone else
		var outboundChangeNameMsg = goWebChat.NewServerChangeNameMessage(oldName, newName)
		outboundChangeNameRaw, err := json.Marshal(outboundChangeNameMsg)

		if err != nil {
			log.Println("Error marshaling message to JSON", outboundChangeNameMsg)
			return
		}

		go BroadcastToAllExcept(newName, outboundChangeNameRaw)
	}
}
