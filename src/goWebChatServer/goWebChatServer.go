package main

import (
	"goWebChat"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

var clientsMap goWebChat.ClientsMap

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
	defer func() { clientsMap.UnregisterClient <- clientPtr }()

	clientsMap.RegisterClient <- clientPtr

	// client process loop
	for {
		select {
		case readBytes := <-client.ReadChan:
			log.Println("On ReadChan: ", string(readBytes))

			handleMessage(readBytes, client)

		case <-client.Closed:
			log.Println("Connection on client ", client.Name, " was closed")
			return
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

	clientsMap.Destroy <- true
}
