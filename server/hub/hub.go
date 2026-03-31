/*
Package hub is blah-blah-blah
*/
package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"patrware/server/models"
	"patrware/structs"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type _Connection struct {
	ID               structs.UUID
	SocketConn       *websocket.Conn
	Endpoint         models.Endpoint
	MessageReadChan  chan structs.Message
	MessageWriteChan chan structs.Message
	isConnected      bool
	pHub             *Hub
}

func newConnection(socketConn *websocket.Conn, endpoint models.Endpoint) *_Connection {
	conn := &_Connection{
		ID:               structs.GenerateUUID(),
		SocketConn:       socketConn,
		Endpoint:         endpoint,
		MessageReadChan:  make(chan structs.Message),
		MessageWriteChan: make(chan structs.Message),
		pHub:             _HubInstance,
	}
	go conn.readHandler()
	go conn.writeHandler()
	return conn
}

func (conn *_Connection) readHandler() {
	for conn.isConnected {
		var message structs.Message
		if err := conn.SocketConn.ReadJSON(&message); err == nil {
			conn.MessageReadChan <- message
		} else {
			log.Println(err.Error())
			conn.isConnected = false
			conn.pHub.removeConnection(conn.ID)
		}
	}
}

func (conn *_Connection) writeHandler() {
	for conn.isConnected {
		select {
		case msg := <-conn.MessageWriteChan:
			if err := conn.SocketConn.WriteJSON(msg); err != nil {
				log.Println(err)
				conn.isConnected = false
				conn.pHub.removeConnection(conn.ID)
			}
		case <-time.After(5 * time.Second):
		}
	}
}

type Hub struct {
	endpoints             map[structs.UUID]*_Connection
	expiredConnectionsIdx map[structs.UUID]struct{}
	connMX                sync.Mutex
	generalMX             sync.Mutex
	upgrader              websocket.Upgrader
}

var _HubInstance *Hub

func InitHub() {
	_HubInstance = &Hub{
		endpoints:             make(map[structs.UUID]*_Connection),
		expiredConnectionsIdx: make(map[structs.UUID]struct{}),
		connMX:                sync.Mutex{},
		generalMX:             sync.Mutex{},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	_HubInstance.setupHandlers()
	go _HubInstance.removeExpiredConnections()
}

func ClearUp() {
	for id, conn := range _HubInstance.endpoints {
		if err := conn.SocketConn.Close(); err == nil {
			_HubInstance.deleteConnection(id)
		} else {
			log.Println(err.Error())
		}
	}
}

func (hub *Hub) setupHandlers() {
	http.HandleFunc("/ws", hub.serveConnection)
}

func (hub *Hub) serveConnection(writer http.ResponseWriter, req *http.Request) {
	conn, err := hub.upgrader.Upgrade(writer, req, nil)
	if err != nil {
		panic("Not implemented")
	}
	var helloMessage structs.Message
	if err = conn.ReadJSON(&helloMessage); err != nil {
		panic(err.Error())
	}
	if helloMessage.Type != structs.MessageTypeHello {
		log.Println("Wrong message type")
		return
	}
	var endpointInfo structs.EndpointInfo
	if err = json.Unmarshal(helloMessage.Payload, &endpointInfo); err != nil {
		log.Println(err.Error())
		return
	}

	endPoint := models.MakeEndpoint(endpointInfo)
	hubconn := newConnection(conn, endPoint)
	hub.addConnection(hubconn)
	// defer func() {
	// 	hub.deleteConnection(hubconn.Id)
	// 	conn.Close()
	// }()
	fmt.Printf("Got a new endpoint:, %v", req.Host)
}

func GetAllEndpoints() []models.Endpoint {
	endpoints := make([]models.Endpoint, len(_HubInstance.endpoints))
	i := 0
	for _, conn := range _HubInstance.endpoints {
		endpoints[i] = conn.Endpoint
		i++
	}
	return endpoints
}

func GetConnectionAssociatedWithEndpoint(endpointID structs.UUID) (*_Connection, error) {
	for _, conn := range _HubInstance.endpoints {
		if conn.Endpoint.GetID().Equals(endpointID) {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("no connection assiciated with the endpoint: %v", endpointID)
}

func (hub *Hub) addConnection(conn *_Connection) {
	hub.generalMX.Lock()
	hub.endpoints[conn.ID] = conn
	hub.generalMX.Unlock()
}

func (hub *Hub) deleteConnection(id structs.UUID) {
	hub.generalMX.Lock()
	delete(hub.endpoints, id)
	hub.generalMX.Unlock()
}

func (hub *Hub) removeExpiredConnections() {
	hub.connMX.Lock()
	defer hub.connMX.Unlock()

	for id := range hub.expiredConnectionsIdx {
		if conn, ok := hub.endpoints[id]; ok {
			if err := conn.SocketConn.Close(); err != nil {
				log.Println(err.Error())
			}
			hub.deleteConnection(id)
		}
	}
	hub.expiredConnectionsIdx = make(map[structs.UUID]struct{})
}

func (hub *Hub) removeConnection(connID structs.UUID) {
	hub.connMX.Lock()
	defer hub.connMX.Unlock()

	hub.expiredConnectionsIdx[connID] = struct{}{}
}
