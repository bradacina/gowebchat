package main

import (
	"encoding/json"
	"goWebChat"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"code.google.com/p/go.net/websocket"
)

var clientsMap goWebChat.ClientsMap

func BroadcastToAll(msg []byte) {
	clients := clientsMap.GetAllClients()

	for _, k := range clients {
		k.WriteChan <- msg
	}
}

func BroadcastToAllExcept(name string, msg []byte) {
	clients := clientsMap.GetAllClients()

	for _, k := range clients {
		if k.Name != name {
			k.WriteChan <- msg
		}
	}
}

func GetUniqueName(name string) string {
	clients := clientsMap.GetAllClients()

	good := true

	for {
		for _, k := range clients {
			if k.Name == name {
				good = false
				break
			}
		}

		if good {
			return name
		} else {
			// attach a random number to the name
			name = name + strconv.Itoa(rand.Intn(10))
			good = true
		}
	}
}

func SendListOfConnectedClients(c *goWebChat.Client) {
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

func ChatHandler(ws *websocket.Conn) {

	req := ws.Request()

	err := req.ParseForm()
	if err != nil {
		ws.Close()
		log.Printf("Can't parse form")
	}

	name, ok := req.Form["name"]

	if !ok {
		log.Printf("Client did not provide a name")
		ws.Close()
		return
	}

	newName := GetUniqueName(name[0])

	client := goWebChat.NewClient(newName, ws)
	clientPtr := &client

	defer log.Println("Exiting handler function.")
	defer clientsMap.UnregisterClient(clientPtr)

	clientsMap.RegisterClient(clientPtr)

	// client process loop
	for {
		select {
		case <-clientPtr.Closed:
			log.Println("Connection on client ", client.Name, " was closed")
			return

		case readBytes := <-clientPtr.ReadChan:
			log.Println("On ReadChan: ", string(readBytes))

			handleMessage(readBytes, client)

		case disconnectedClient := <-clientsMap.ClientUnregistered:

			// notify everyone that a user has disconnected
			var outboundChatMsg = goWebChat.NewServerClientPartMessage(disconnectedClient.Name)
			outboundRaw, err := json.Marshal(outboundChatMsg)

			if err != nil {
				log.Println("Error marshaling message to JSON", outboundChatMsg)
				return
			}

			go BroadcastToAll(outboundRaw)

		case connectedClient := <-clientsMap.ClientRegistered:
			// notify everyone that a new user has connected
			var outboundChatMsg = goWebChat.NewServerClientJoinMessage(connectedClient.Name)
			outboundRaw, err := json.Marshal(outboundChatMsg)

			if err != nil {
				log.Println("Error marshaling message to JSON", outboundChatMsg)
				return
			}

			go BroadcastToAllExcept(connectedClient.Name, outboundRaw)

			go SendListOfConnectedClients(connectedClient)

		}
	}
}

func main() {

	clientsMap = goWebChat.NewClientsMap()

	http.Handle("/", http.FileServer(http.Dir("../../html")))
	http.Handle("/chat", websocket.Handler(ChatHandler))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
