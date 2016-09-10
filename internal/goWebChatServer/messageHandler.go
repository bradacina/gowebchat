package goWebChatServer

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/bradacina/gowebchat/internal/goWebChat"
)

func handleCommand(msg string, client *goWebChat.Client) {
	log.Println("Got command message: ", msg)

	tokens := strings.Split(msg, " ")
	if len(tokens) == 0 {
		return
	}

	cmd := tokens[0]
	args := tokens[1:]

	switch cmd {
	case "login":
		if len(args) == 0 {
			return
		}

		if args[0] == "password987" {
			client.SetIsAdmin(true)
			sendStatusMessage(client, "You are now Admin.")
		}
	case "logout":
		client.SetIsAdmin(false)
		sendStatusMessage(client, "You are no longer Admin.")
	case "whois":
		if len(args) == 0 || !client.IsAdmin() {
			return
		}

		whoisClient, ok := clientsMap.GetClient(args[0])

		if ok {
			msg := fmt.Sprintf("<b>%v</b> (Admin:<b>%v</b>) has connected from <b>%v</b> using\n<b>%v</b>.",
				whoisClient.Name, whoisClient.IsAdmin, whoisClient.IPAddr, whoisClient.UserAgent)
			sendStatusMessage(client, msg)
		}

	case "kick":
		if len(args) == 0 || !client.IsAdmin() {
			return
		}

		clientToKick, ok := clientsMap.GetClient(args[0])

		if ok {
			clientToKick.Close()
		}
	}
}

func handleMessage(msg []byte, client *goWebChat.Client) {

	var messageType goWebChat.MessageType

	json.Unmarshal(msg, &messageType)
	log.Println("MessageType ", messageType)

	// did we receive a malformed message
	if messageType == (goWebChat.MessageType{}) {
		return
	}

	switch messageType.Type {
	case "ClientPong":
		pongMsg, err := goWebChat.UnmarshalClientPongMessage(msg)

		if err != nil {
			log.Println("Could not unmarshal incoming message from client: ", msg)
			return
		}

		// if we got the correct pong reply then stop the ping timer
		if client.PingPayload() == pongMsg.Payload {
			client.ResetPingTimeout()
		}

	case "ClientChat":
		chatMsg, err := goWebChat.UnmarshalClientChatMessage(msg)

		if err != nil || len(chatMsg.Chat) == 0 {
			log.Println("Could not unmarshall incoming message from client: ", msg)
			return
		}

		// treat special case if chat message starts with / which means it's a command
		if chatMsg.Chat[0] == '/' {
			handleCommand(chatMsg.Chat[1:], client)
			return
		}

		chatMessageContent := html.EscapeString(chatMsg.Chat)
		var outboundChatMsg = goWebChat.NewServerChatMessage(chatMessageContent, client.Name())

		go broadcastToAll(outboundChatMsg)
	case "ChangeName":
		changeNameMsg, err := goWebChat.UnmarshalClientChangeNameMessage(msg)

		if err != nil || len(changeNameMsg.NewName) == 0 {
			log.Println("Could not unmarshall incoming message from client: ", msg)
			return
		}

		if client.Name() == changeNameMsg.NewName {
			return
		}

		// set new name
		oldName := string([]byte(client.Name()))
		newName := changeName(client.Name(), changeNameMsg.NewName)

		log.Println(oldName, newName, changeNameMsg.NewName)

		sendName(newName, client)

		// broadcast the name change to everyone else
		var outboundChangeNameMsg = goWebChat.NewServerChangeNameMessage(oldName, newName)
		go broadcastToAllExcept(newName, outboundChangeNameMsg)
	}
}
