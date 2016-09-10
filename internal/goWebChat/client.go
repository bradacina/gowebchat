package goWebChat

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

// Client stores internal info about the connected client.
type Client struct {
	ReadChan chan []byte
	Closed   chan interface{}

	name        string
	ipAddr      string
	userAgent   string
	pingPayload int
	isAdmin     bool
	pingTimeout <-chan time.Time

	con  *websocket.Conn
	lock sync.RWMutex
}

// NewClient creates a new client type given a name, websocket connection,
// user agent and ip address.
func NewClient(name string, ws *websocket.Conn, userAgent string, ipAddress string) *Client {
	client := Client{name: name, con: ws}
	client.ReadChan = make(chan []byte, 10)
	client.Closed = make(chan interface{})
	client.ipAddr = ipAddress
	client.userAgent = userAgent

	go client.readLoop()

	return &client
}

// Name returns the name used by the client in the chat.
func (c *Client) Name() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.name
}

// SetName sets the name used by the client.
func (c *Client) SetName(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.name = name
}

// IPAddr returns the ip address of the client.
func (c *Client) IPAddr() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.ipAddr
}

// UserAgent returns the user agent of the client's browser.
func (c *Client) UserAgent() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.userAgent
}

// Close disconnects a client from the server.
func (c *Client) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	err := c.con.Close()
	c.con = nil
	return err
}

func (c *Client) logMessage(a ...interface{}) {
	var newMsg []interface{}
	newMsg = append(newMsg, c.name+"->")
	newMsg = append(newMsg, a...)
	log.Println(newMsg...)
}

func (c *Client) readLoop() {

	for {
		bytes := make([]byte, 1024)
		nBytes, err := c.con.Read(bytes)
		if err != nil {
			c.lock.RLock()
			c.logMessage("Got error when reading from ", c.name)
			c.lock.RUnlock()
			break
		}

		if nBytes == 0 {
			c.logMessage("Read 0 bytes")
		}

		// send the data down the channel
		c.ReadChan <- bytes[:nBytes]

		//c.logMessage("Read ", bytes[0:nBytes])
	}

	c.lock.Lock()
	c.con.Close()
	c.con = nil
	c.lock.Unlock()
	c.Closed <- true

	c.logMessage("Read loop exiting.")
}

// Send sends a byte slice to the client.
func (c *Client) Send(msg []byte) {
	if c.con != nil {
		c.con.Write(msg)
	}
}

// IsAdmin returns true if the client has admin privileges.
func (c *Client) IsAdmin() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.isAdmin
}

// SetIsAdmin will give/remove to/from a client admin privileges.
func (c *Client) SetIsAdmin(admin bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.isAdmin = admin
}

// PingPayload returns the payload sent to the client with the last ping.
func (c *Client) PingPayload() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.pingPayload
}

// SetPingPayload saves the payload sent to the client with the last ping.
func (c *Client) SetPingPayload(payload int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.pingPayload = payload
}

// PingTimeout returns the pingTimeout channel.
func (c *Client) PingTimeout() <-chan time.Time {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.pingTimeout
}

// SetPingTimeout sets the pingTimeout channel for the client.
func (c *Client) SetPingTimeout(timeout <-chan time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.pingTimeout = timeout
}

// ResetPingTimeout sets the pingTimeout channel to nil.
func (c *Client) ResetPingTimeout() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.pingTimeout = nil
}
