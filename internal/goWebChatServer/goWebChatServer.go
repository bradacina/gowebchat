package goWebChatServer

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/bradacina/gowebchat/internal/goWebChat"

	"golang.org/x/net/websocket"
)

// Server represents a web chat server.``
type Server struct {
	clientsMap *goWebChat.ClientsMap
}

// NewServer creates a new web chat server instance
func NewServer() *Server {
	s := Server{}
	s.clientsMap = goWebChat.NewClientsMap()
	return &s
}

// ChatHandler handles an incoming websocket connection.
func (s *Server) ChatHandler(ws *websocket.Conn) {

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

	newName := s.getUniqueName(name[0])

	client := goWebChat.NewClient(newName, ws, req.UserAgent(), getIP(req))

	defer log.Println("Exiting handler function.")
	defer s.unregisterClient(client)

	s.clientsMap.RegisterClient(client)

	// notify everyone that a new user has connected
	var outboundChatMsg = goWebChat.NewServerClientJoinMessage(client.Name())

	go s.broadcastToAllExcept(client.Name(), outboundChatMsg)

	go sendName(client.Name(), client)

	go s.sendListOfConnectedClients(client)

	// client process loop
	for {

		activityTimeout := time.After(30 * time.Second)

		select {
		case <-client.Closed:
			log.Println("Connection on client ", client.Name(), " was closed")
			return

		case readBytes := <-client.ReadChan:
			log.Println("On ReadChan: ", string(readBytes))

			s.handleMessage(readBytes, client)

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

func (s *Server) getUniqueName(name string) string {
	clients := s.clientsMap.GetAllClients()

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

func (s *Server) changeName(oldName string, newName string) string {
	uniqueName := s.getUniqueName(cleanupName(newName))

	s.clientsMap.ReplaceName(oldName, uniqueName)

	return uniqueName
}

func (s *Server) unregisterClient(client *goWebChat.Client) {
	if err := s.clientsMap.UnregisterClient(client); err != nil {
		return
	}

	// notify everyone that a user has disconnected (except the disconnected user of course)
	var outboundChatMsg = goWebChat.NewServerClientPartMessage(client.Name())

	go s.broadcastToAllExcept(client.Name(), outboundChatMsg)
}
