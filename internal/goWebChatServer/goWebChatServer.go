package goWebChatServer

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bradacina/gowebchat/internal/goWebChat"

	"golang.org/x/net/websocket"
)

var clientsMap *goWebChat.ClientsMap

func getUniqueName(name string) string {
	clients := clientsMap.GetAllClients()

	isGood := true

	for {
		for _, k := range clients {
			if k.Name() == name {
				isGood = false
				break
			}
		}

		if isGood {
			log.Println("Generated unique name: ", name)
			return name
		}

		// attach a random number to the name
		name = name + strconv.Itoa(rand.Intn(10))
		isGood = true
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

func getIP(req *http.Request) string {

	ipSlice, ok := req.Header["X-Real-Ip"]
	var ip string
	if !ok {
		ip = req.RemoteAddr
	} else {
		ip = ipSlice[0]
	}
	index := strings.Index(ip, ":")
	if index != -1 {
		ip = ip[0:index]
	}

	return ip
}

// ChatHandler handles an incoming websocket connection.
func ChatHandler(ws *websocket.Conn) {

	log.Println("In ChatHandler()")
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

	client := goWebChat.NewClient(newName, ws, req.UserAgent(), getIP(req))

	defer log.Println("Exiting handler function.")
	defer unregisterClient(client)

	clientsMap.RegisterClient(client)

	// notify everyone that a new user has connected
	var outboundChatMsg = goWebChat.NewServerClientJoinMessage(client.Name())

	go broadcastToAllExcept(client.Name(), outboundChatMsg)

	go sendName(client.Name(), client)

	go sendListOfConnectedClients(client)

	// client process loop
	for {

		activityTimeout := time.After(30 * time.Second)

		select {
		case <-client.Closed:
			log.Println("Connection on client ", client.Name(), " was closed")
			return

		case readBytes := <-client.ReadChan:
			log.Println("On ReadChan: ", string(readBytes))

			handleMessage(readBytes, client)

		case <-activityTimeout:
			// send ping
			if client.PingTimeout() == nil {
				go sendPing(client)
				client.SetPingTimeout(time.After(30 * time.Second))
			}

		case <-client.PingTimeout():
			// we did not receive a reply to our timeout so disconnect this user
			client.Close()
			log.Println("User ", client.Name(), " timed out.")
		}
	}
}

func unregisterClient(client *goWebChat.Client) {
	if err := clientsMap.UnregisterClient(client); err != nil {
		return
	}

	// notify everyone that a user has disconnected (except the disconnected user of course)
	var outboundChatMsg = goWebChat.NewServerClientPartMessage(client.Name())

	go broadcastToAllExcept(client.Name(), outboundChatMsg)
}

func init() {
	clientsMap = goWebChat.NewClientsMap()
}
