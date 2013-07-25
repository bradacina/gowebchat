package main

import (
	"encoding/json"
	"goWebChat"
	"log"
)

func sendStatusMessage(client *goWebChat.Client, text string) {

	var outboundChatMsg = goWebChat.NewServerStatusMessage(text)
	outboundRaw, err := json.Marshal(outboundChatMsg)

	if err != nil {
		log.Println("Error marshaling message to JSON", outboundChatMsg)
		return
	}

	client.WriteChan <- outboundRaw
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
		k.WriteChan <- msg
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
		if k.Name != name {
			k.WriteChan <- msg
		}
	}
}

func sendListOfConnectedClients(c *goWebChat.Client) {
	var users string
	clients := clientsMap.GetAllClients()

	for _, k := range clients {
		users = users + "," + k.Name
	}

	var outboundClientListMsg = goWebChat.NewServerClientListMessage(users)
	outboundRaw, err := json.Marshal(outboundClientListMsg)

	if err != nil {
		log.Println("Error marshaling message to JSON", outboundClientListMsg)
		return
	}

	c.WriteChan <- outboundRaw
}

func sendName(name string, client *goWebChat.Client) {
	// send the new name back to client
	var outboundSetNameMsg = goWebChat.NewServerSetNameMessage(name)
	outboundSetNameRaw, err := json.Marshal(outboundSetNameMsg)

	if err != nil {
		log.Println("Error marshaling message to JSON", outboundSetNameMsg)
		return
	}

	client.WriteChan <- outboundSetNameRaw
}

func sendPing(client *goWebChat.Client) {
	// send the ping message
	client.PingPayload = random(123524, 967845)
	var outboundPingMsg = goWebChat.NewServerPingMessage(client.PingPayload)
	outboundPingRaw, err := json.Marshal(outboundPingMsg)

	if err != nil {
		log.Println("Error marshalling message to JSON", outboundPingMsg)
		return
	}

	client.WriteChan <- outboundPingRaw
}
