package main

import (
	"goWebChat"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

var clients []*goWebChat.Client
var registerClient chan *goWebChat.Client
var unregisterClient chan *goWebChat.Client

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

	client := goWebChat.NewClient(name[0], ws)
	clientPtr := &client

	defer client.Close()
	defer func() { unregisterClient <- clientPtr }()

	log.Println("Number of clients before: ", len(clients))

	registerClient <- clientPtr

	// client process loop
	for {
		select {
		case readBytes := <-client.ReadChan:
			log.Println("On ReadChan: ", readBytes)

			go broadcastToAllClients(append([]byte(client.Name+" said-> "), readBytes...))

		case <-client.Closed:
			log.Println("Connection on client ", client.Name, " was closed")
			return
		}
	}
}

func broadcastToAllClients(msg []byte) {
	for _, val := range clients {
		val.WriteChan <- msg
	}
}

func registrationLoop() {
	for {
		select {
		case addClient := <-registerClient:
			clients = append(clients, addClient)
			log.Println("Adding client to array: ", addClient.Name)

			broadcastToAllClients([]byte(addClient.Name + " has connected"))
		case removeClient := <-unregisterClient:

			log.Println("Removing client from array: ", removeClient.Name)
			index := -1
			for k, v := range clients {
				if v == removeClient {
					index = k
					break
				}
			}

			if index != -1 {
				clients = append(clients[:index], clients[index+1:]...)
				broadcastToAllClients([]byte(removeClient.Name + " has disconnected"))
			}
		}
	}
}

func main() {

	clients = make([]*goWebChat.Client, 0)
	registerClient = make(chan *goWebChat.Client, 10)
	unregisterClient = make(chan *goWebChat.Client, 10)

	go registrationLoop()

	http.Handle("/", http.FileServer(http.Dir("../../html")))
	http.Handle("/chat", websocket.Handler(ChatHandler))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
