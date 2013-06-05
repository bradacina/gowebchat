package goWebChat

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

type Client struct {
	Name      string
	ReadChan  chan []byte
	WriteChan chan []byte
	Closed    chan bool

	con        *websocket.Conn
	writeClose chan bool
}

func (c *Client) Close() error {
	return c.con.Close()
}

func NewClient(name string, ws *websocket.Conn) Client {

	client := Client{Name: name, con: ws}
	client.ReadChan = make(chan []byte, 10)
	client.WriteChan = make(chan []byte, 10)
	client.Closed = make(chan bool, 0)
	client.writeClose = make(chan bool, 0)

	go client.readLoop()
	go client.writeLoop()

	return client
}

func (c *Client) logMessage(a ...interface{}) {
	newMsg := make([]interface{}, 0)
	newMsg = append(newMsg, c.Name+"->")
	newMsg = append(newMsg, a...)
	log.Println(newMsg...)
}

func (c *Client) readLoop() {

	for {
		bytes := make([]byte, 1024)
		nBytes, err := c.con.Read(bytes)
		if err != nil {
			c.logMessage("Got error when reading from ", c.Name)
			break
		}

		if nBytes == 0 {
			c.logMessage("Read 0 bytes")
		}

		// send the data down the channel
		copyBytes := make([]byte, nBytes)
		copy(copyBytes, bytes[0:nBytes])
		c.ReadChan <- copyBytes

		c.logMessage("Read ", bytes[0:nBytes])
	}
	c.con.Close()
	c.Closed <- true
	c.writeClose <- true

	c.logMessage("Read loop exiting.")
}

func (c *Client) writeLoop() {
	for {
		select {
		case writeBytes := <-c.WriteChan:
			c.logMessage("On WriteChan: ", writeBytes)
			c.con.Write(writeBytes)
		case <-c.writeClose:
			c.logMessage("Write Loop stopping due to client closed event")
			return
		}
	}
}
