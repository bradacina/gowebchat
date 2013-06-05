package goWebChat

import (
	"log"
)

type ClientsMap struct {
	RegisterClient       chan *Client
	UnregisterClient     chan *Client
	BroadcastToAll       chan []byte
	BroadcastToAllExcept chan BroadcastToAllExcept
	Destroy              chan bool
	GetAllClients        chan chan []*Client
	GetUniqueName        chan string

	ClientRegistered   chan *Client
	ClientUnregistered chan *Client

	// todo: replace the slice with a map
	clients []*Client
}

type BroadcastToAllExcept struct {
	Name    string
	Content []byte
}

func NewClientsMap() ClientsMap {

	clientsMap := ClientsMap{}

	clientsMap.clients = make([]*Client, 0)
	clientsMap.RegisterClient = make(chan *Client, 10)
	clientsMap.UnregisterClient = make(chan *Client, 10)
	clientsMap.BroadcastToAll = make(chan []byte, 10)
	clientsMap.BroadcastToAllExcept = make(chan BroadcastToAllExcept, 10)
	clientsMap.Destroy = make(chan bool, 0)
	clientsMap.GetAllClients = make(chan chan []*Client, 10)
	clientsMap.GetUniqueName = make(chan chan string)
	clientsMap.ClientRegistered = make(chan *Client, 10)
	clientsMap.ClientUnregistered = make(chan *Client, 10)

	go clientsMap.loop()

	return clientsMap
}

func (cM *ClientsMap) loop() {
	for {
		select {
		case addClient := <-cM.RegisterClient:
			cM.clients = append(cM.clients, addClient)
			log.Println("Added client to array: ", addClient.Name)
			log.Println("Number of clients currently in map: ", len(cM.clients))

			cM.ClientRegistered <- addClient

		case removeClient := <-cM.UnregisterClient:
			index := -1
			for k, v := range cM.clients {
				if v == removeClient {
					index = k
					break
				}
			}

			if index != -1 {
				cM.clients = append(cM.clients[:index], cM.clients[index+1:]...)

				cM.ClientUnregistered <- removeClient
			}

			log.Println("Removed client from array: ", removeClient.Name)
			log.Println("Current number of clients: ", len(cM.clients))

		case msg := <-cM.BroadcastToAll:
			for _, val := range cM.clients {
				val.WriteChan <- msg
			}

		case msg := <-cM.BroadcastToAllExcept:
			for _, val := range cM.clients {
				if val.Name != msg.Name {
					val.WriteChan <- msg.Content
				}
			}

		case returnChan := <-cM.GetAllClients:
			clients := make([]*Client, len(cM.clients), len(cM.clients))
			copy(clients, cM.clients)
			returnChan <- clients

		case <-cM.Destroy:
			log.Println("Destroying clients map")

		}
	}
}
