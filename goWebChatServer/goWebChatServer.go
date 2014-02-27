package main

import (
	"code.metaconstudios.com/gowebchat/goWebChat"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/go.net/websocket"
)

var clientsMap goWebChat.ClientsMap

func getUniqueName(name string) string {
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

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func changeName(oldName string, newName string) string {
	uniqueName := getUniqueName(cleanupName(newName))

	clientsMap.ReplaceName(oldName, uniqueName)

	return uniqueName
}

func cleanupName(oldName string) string {
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

func getIp(req *http.Request) string {
	ipslice, ok := req.Header["X-Real-Ip"]
	var ip string
	if !ok {
		ip = req.RemoteAddr
	} else {
		ip = ipslice[0]
	}
	index := strings.Index(ip, ":")
	if index != -1 {
		ip = ip[0:index]
	}

	return ip
}

func chatHandler(ws *websocket.Conn) {

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

	newName := getUniqueName(name[0])

	client := goWebChat.NewClient(newName, ws, req.UserAgent(), getIp(req))
	clientPtr := &client

	defer log.Println("Exiting handler function.")
	defer clientsMap.UnregisterClient(clientPtr)

	clientsMap.RegisterClient(clientPtr)

	// client process loop
	for {

		clientPtr.ActivityTimeout = time.After(30 * time.Second)

		select {
		case <-clientPtr.Closed:
			log.Println("Connection on client ", client.Name, " was closed")
			return

		case readBytes := <-clientPtr.ReadChan:
			log.Println("On ReadChan: ", string(readBytes))

			handleMessage(readBytes, clientPtr)

		case <-clientPtr.ActivityTimeout:
			// send ping
			if clientPtr.PingTimeout == nil {
				go sendPing(clientPtr)
				clientPtr.PingTimeout = time.After(30 * time.Second)
			}

		case <-clientPtr.PingTimeout:
			// we did not receive a reply to our timeout so disconnect this user
			clientPtr.Close()
			log.Println("User ", client.Name, " timed out.")

		}
	}
}

func clientMapLoop() {
	for {
		select {
		case disconnectedClient := <-clientsMap.ClientUnregistered:
			// notify everyone that a user has disconnected (except the disconnected user of course)
			var outboundChatMsg = goWebChat.NewServerClientPartMessage(disconnectedClient.Name)

			go broadcastToAllExcept(disconnectedClient.Name, outboundChatMsg)

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
	go clientMapLoop()

	http.Handle("/", http.FileServer(http.Dir("../../html")))
	http.Handle("/chat", websocket.Handler(chatHandler))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
