package main

import (
	"net/http"

	"github.com/bradacina/gowebchat/internal/goWebChatServer"
	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("html")))
	http.Handle("/chat", websocket.Handler(goWebChatServer.ChatHandler))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
