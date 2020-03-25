// +build !appengine

package main

import (
        "net/http"

        "github.com/bradacina/gowebchat/internal/goWebChatServer"
        "golang.org/x/net/websocket"
)

var server *goWebChatServer.Server

func init() {
        server := goWebChatServer.NewServer()
        http.Handle("/", http.FileServer(http.Dir("html")))
        http.Handle("/chat", websocket.Handler(server.ChatHandler))
}


func main() {
	err := http.ListenAndServe(":54321", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
