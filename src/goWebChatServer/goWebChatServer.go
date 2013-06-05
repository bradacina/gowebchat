package main

import (
	"encoding/json"
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

	uniqueName := make(chan string)
	var newName string

	clientsMap.GetUniqueName <- name[0]
	select {
	case newName = <-uniqueName:

	}

	client := goWebChat.NewClient(newname, ws)
	clientPtr := &client

	defer log.Println("Exiting handler function.")
	//defer client.Close()
	defer func() { clientsMap.UnregisterClient <- clientPtr }()

	clientsMap.RegisterClient <- clientPtr

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

			clientsMap.BroadcastToAll <- outboundRaw

		case connectedClient := <-clientsMap.ClientRegistered:
			// notify everyone that a new user has connected
			var outboundChatMsg = goWebChat.NewServerClientJoinMessage(connectedClient.Name)
			outboundRaw, err := json.Marshal(outboundChatMsg)

			if err != nil {
				log.Println("Error marshaling message to JSON", outboundChatMsg)
				return
			}

			clientsMap.BroadcastToAllExcept <- goWebChat.BroadcastToAllExcept{Name: connectedClient.Name, Content: outboundRaw}

			retChan := make(chan []*goWebChat.Client)

			clientsMap.GetAllClients <- retChan

			// send the list of users to the newly connected user
			go func() {
				var users string

				select {
				case clients := <-retChan:

					for _, k := range clients {
						users = users + "," + k.Name
					}
				}

				var outboundClientListMsg = goWebChat.NewServerClientListMessage(users)
				outboundRaw, err = json.Marshal(outboundClientListMsg)

				if err != nil {
					log.Println("Error marshaling message to JSON", outboundChatMsg)
					return
				}

				connectedClient.WriteChan <- outboundRaw
			}()

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
