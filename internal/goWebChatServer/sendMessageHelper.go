package goWebChatServer

import (
	"encoding/json"
	"log"

	"github.com/bradacina/gowebchat/internal/goWebChat"
)

func sendStatusMessage(client *goWebChat.Client, text string) {

	var outboundChatMsg = goWebChat.NewServerStatusMessage(text)
	outboundRaw, err := json.Marshal(outboundChatMsg)

	if err != nil {
		log.Println("Error marshaling message to JSON", outboundChatMsg)
		return
	}

	client.Send(outboundRaw)
}

func broadcastToAll(outbound interface{}) {
	defer log.Println("Exiting BroadcastToAll")
	log.Println("Entered BroadcastToAll")

	msg, err := json.Marshal(outbound)

	if err != nil {
		log.Println("Error marshaling message to JSON", outbound)
		return
	}

	clients := clientsMap.GetAllClients()

	for _, k := range clients {
		k.Send(msg)
	}
}

func broadcastToAllExcept(name string, outbound interface{}) {

	defer log.Println("Exiting BroadcastToAllExcept")
	log.Println("Entering BroadcastToAllExcept")

	msg, err := json.Marshal(outbound)

	if err != nil {
		log.Println("Error marshaling message to JSON", outbound)
		return
	}

	clients := clientsMap.GetAllClients()

	for _, k := range clients {
		if k.Name() != name {
			k.Send(msg)
		}
	}
}

func sendListOfConnectedClients(c *goWebChat.Client) {
	var users string
	clients := clientsMap.GetAllClients()

	for _, k := range clients {
		users = users + "," + k.Name()
	}

	var outboundClientListMsg = goWebChat.NewServerClientListMessage(users)
	outboundRaw, err := json.Marshal(outboundClientListMsg)

	if err != nil {
		log.Println("Error marshaling message to JSON", outboundClientListMsg)
		return
	}

	c.Send(outboundRaw)
}

func sendName(name string, client *goWebChat.Client) {
	// send the new name back to client
	var outboundSetNameMsg = goWebChat.NewServerSetNameMessage(name)
	outboundSetNameRaw, err := json.Marshal(outboundSetNameMsg)

	if err != nil {
		log.Println("Error marshaling message to JSON", outboundSetNameMsg)
		return
	}

	client.Send(outboundSetNameRaw)
}

func sendPing(client *goWebChat.Client) {
	// send the ping message
	pingPayload := random(123524, 967845)
	client.SetPingPayload(pingPayload)
	var outboundPingMsg = goWebChat.NewServerPingMessage(pingPayload)
	outboundPingRaw, err := json.Marshal(outboundPingMsg)

	if err != nil {
		log.Println("Error marshalling message to JSON", outboundPingMsg)
		return
	}

	client.Send(outboundPingRaw)
}
