// the map keeps track of connected clients

package goWebChat

import (
	"errors"
	"log"
	"sync"
)

type ClientsMap struct {
	ClientRegistered   chan *Client
	ClientUnregistered chan *Client

	// todo: replace the slice with a map
	clients []*Client

	lock sync.Mutex
}

func NewClientsMap() ClientsMap {

	clientsMap := ClientsMap{}
	clientsMap.ClientRegistered = make(chan *Client, 10)
	clientsMap.ClientUnregistered = make(chan *Client, 10)

	clientsMap.clients = make([]*Client, 0)

	return clientsMap
}

func (cM *ClientsMap) RegisterClient(c *Client) {
	cM.lock.Lock()
	defer cM.lock.Unlock()
	cM.clients = append(cM.clients, c)

	log.Println("Added client to map. New size of map: ", len(cM.clients))
	cM.ClientRegistered <- c
}

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
		cM.ClientUnregistered <- c
		return nil
	} else {
		log.Println("Could not find client in map to remove.")
		return errors.New("Could not find client to remove.")
	}
}

func (cM *ClientsMap) GetAllClients() []*Client {
	cM.lock.Lock()
	defer cM.lock.Unlock()

	result := make([]*Client, len(cM.clients))
	copy(result, cM.clients)
	return result
}

func (cM *ClientsMap) ReplaceName(toReplace string, replacement string) {
	cM.lock.Lock()
	defer cM.lock.Unlock()

	for i, k := range cM.clients {
		if k.Name == toReplace {
			cM.clients[i].Name = replacement
		}
	}
}

func (cM *ClientsMap) GetClient(name string) (*Client, bool) {
	cM.lock.Lock()
	defer cM.lock.Unlock()

	for _, k := range cM.clients {
		if k.Name == name {
			return k, true
		}
	}
	return nil, false
}
