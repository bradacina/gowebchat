// the map keeps track of connected clients

package goWebChat

import (
	"errors"
	"log"
	"sync"
)

// ClientsMap keeps track of what clients are currently connected
type ClientsMap struct {
	// todo: replace the slice with a map
	clients []*Client

	lock sync.Mutex
}

// NewClientsMap creates and returns a new map that keeps track of clients.
func NewClientsMap() *ClientsMap {

	clientsMap := ClientsMap{}

	clientsMap.clients = make([]*Client, 0)

	return &clientsMap
}

// RegisterClient adds a newly connected client to the map.
func (cM *ClientsMap) RegisterClient(c *Client) {
	cM.lock.Lock()
	defer cM.lock.Unlock()
	cM.clients = append(cM.clients, c)

	log.Println("Added client to map. New size of map: ", len(cM.clients))
}

// UnregisterClient removes a disconnected client from the map.
func (cM *ClientsMap) UnregisterClient(c *Client) error {
	cM.lock.Lock()
	defer cM.lock.Unlock()

	index := -1
	for k, v := range cM.clients {
		if v == c {
			index = k
			break
		}
	}

	if index != -1 {
		cM.clients[index] = nil
		cM.clients = append(cM.clients[:index], cM.clients[index+1:]...)
		log.Println("Removed client from map. New size of map: ", len(cM.clients))
		return nil
	}

	log.Println("Could not find client in map to remove.")
	return errors.New("Could not find client to remove.")
}

// GetAllClients returns a copy of the list of clients that are currently connected.
func (cM *ClientsMap) GetAllClients() []*Client {
	cM.lock.Lock()
	defer cM.lock.Unlock()

	result := make([]*Client, len(cM.clients))
	copy(result, cM.clients)
	return result
}

// ReplaceName changes the name used by a client with a replacement name.
func (cM *ClientsMap) ReplaceName(toReplace string, replacement string) {
	cM.lock.Lock()
	defer cM.lock.Unlock()

	for i, k := range cM.clients {
		if k.Name() == toReplace {
			cM.clients[i].SetName(replacement)
		}
	}
}

// GetClient returns a client by the name he/she uses.
func (cM *ClientsMap) GetClient(name string) (*Client, bool) {
	cM.lock.Lock()
	defer cM.lock.Unlock()

	for _, k := range cM.clients {
		if k.Name() == name {
			return k, true
		}
	}
	return nil, false
}
