package main

import (
	"goWebChat"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"code.google.com/p/go.net/websocket"
)

var clientsMap goWebChat.ClientsMap

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
			log.Println("Generated unique name: ", name)
			return name
		} else {
			// attach a random number to the name
			name = name + strconv.Itoa(rand.Intn(10))
			good = true
		}
	}
}

func ChangeName(oldName string, newName string) string {
	uniqueName := GetUniqueName(CleanupName(newName))

	clientsMap.ReplaceName(oldName, uniqueName)

	return uniqueName
}

func CleanupName(oldName string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-',
			r == '_':
			return r
		}
		return -1
	}, oldName)
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

	client := goWebChat.NewClient(newName, ws, req.UserAgent(), req.RemoteAddr)
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

			handleMessage(readBytes, clientPtr)

		case disconnectedClient := <-clientsMap.ClientUnregistered:

			// notify everyone that a user has disconnected
			var outboundChatMsg = goWebChat.NewServerClientPartMessage(disconnectedClient.Name)

			go broadcastToAll(outboundChatMsg)

		case connectedClient := <-clientsMap.ClientRegistered:
			// notify everyone that a new user has connected
			var outboundChatMsg = goWebChat.NewServerClientJoinMessage(connectedClient.Name)

			go broadcastToAllExcept(connectedClient.Name, outboundChatMsg)

			go sendName(connectedClient.Name, connectedClient)

			go sendListOfConnectedClients(connectedClient)
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
