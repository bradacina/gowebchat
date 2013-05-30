package main

import (
	"io"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

func ChatHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func main() {
	http.Handle("/chat", websocket.Handler(ChatHandler))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
