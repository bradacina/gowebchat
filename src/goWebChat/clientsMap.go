package goWebChat

import (
	"encoding/json"
	"log"
)

type ClientsMap struct {
	RegisterClient   chan *Client
	UnregisterClient chan *Client
	BroadcastToAll   chan []byte
	Destroy          chan bool

	// todo: replace the slice with a map
	clients []*Client
}

func NewClientsMap() ClientsMap {

	clientsMap := ClientsMap{}

	clientsMap.clients = make([]*Client, 0)
	clientsMap.RegisterClient = make(chan *Client, 10)
	clientsMap.UnregisterClient = make(chan *Client, 10)
	clientsMap.BroadcastToAll = make(chan []byte, 10)
	clientsMap.Destroy = make(chan bool, 0)

	go clientsMap.loop()

	return clientsMap
}

func (cM *ClientsMap) loop() {
	for {
		select {
		case addClient := <-cM.RegisterClient:
			cM.clients = append(cM.clients, addClient)
			log.Println("Adding client to array: ", addClient.Name)
			log.Println("Number of clients currently in map: ", len(cM.clients))

			// notify everyone that a new user has connected
			var outboundChatMsg = NewServerStatusMessage(addClient.Name + " has connected")
			outboundRaw, err := json.Marshal(outboundChatMsg)

			if err != nil {
				log.Println("Error marshaling message to JSON", outboundChatMsg)
				return
			}

			cM.BroadcastToAll <- outboundRaw

			// send the list of users to the newly connected user
			var users string
			for i, k := range cM.clients {
				users = users + "," + k
			}
			var outboundClientListMsg = NewServerClientListMessage(users)
			ountboundRaw, err := json.Marshal(outboundClientListMsg)

			if err != nil {
				log.Println("Error marshaling message to JSON", outboundChatMsg)
				return
			}

			addClient.WriteChan <- outboundRaw

		case removeClient := <-cM.UnregisterClient:
			log.Println("Removing client from array: ", removeClient.Name)
			index := -1
			for k, v := range cM.clients {
				if v == removeClient {
					index = k
					break
				}
			}

			if index != -1 {
				cM.clients = append(cM.clients[:index], cM.clients[index+1:]...)

				var outboundChatMsg = NewServerStatusMessage(removeClient.Name + " has disconnected")
				outboundRaw, err := json.Marshal(outboundChatMsg)

				if err != nil {
					log.Println("Error marshaling message to JSON", outboundChatMsg)
					return
				}

				cM.BroadcastToAll <- outboundRaw
			}

		case msg := <-cM.BroadcastToAll:
			for _, val := range cM.clients {
				val.WriteChan <- msg
			}

		case <-cM.Destroy:
			log.Println("Destroying clients map")
		}
	}
}
