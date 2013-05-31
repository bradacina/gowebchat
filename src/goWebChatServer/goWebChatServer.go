package main

import (
	"io"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

type Client struct {
	Name string
	Con  *websocket.Conn
}

var clients []Client

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
	}

	client := Client{Name: name[0], Con: ws}

	clients = append(clients, client)

	log.Print(len(clients))

	io.Copy(ws, ws)
}

func main() {

	clients = make([]Client, 0)

	http.Handle("/", http.FileServer(http.Dir("../../html")))
	http.Handle("/chat", websocket.Handler(ChatHandler))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
