package hub

import (
	"fmt"
	"net/http"
	"patrware/server/models"
	"patrware/structs"
	"sync"

	"github.com/gorilla/websocket"
)

type _Connection struct {
	Id         structs.UUID
	SocketConn *websocket.Conn
	Endpoint   models.Endpoint
}

func newConnection(socketConn *websocket.Conn, endpoint models.Endpoint) *_Connection {
	return &_Connection{
		Id:         structs.GenerateUUID(),
		SocketConn: socketConn,
		Endpoint:   endpoint,
	}
}

type Hub struct {
	endpoints map[structs.UUID]*_Connection
	mx        sync.Mutex
	upgrader  websocket.Upgrader
}

func NewHub() Hub {
	return Hub{
		endpoints: make(map[structs.UUID]*_Connection),
		mx:        sync.Mutex{},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (hub *Hub) ServeConnection(writer http.ResponseWriter, req *http.Request) {
	conn, err := hub.upgrader.Upgrade(writer, req, nil)
	if err != nil {
		panic("Not implemented")
	}
	var endpointInfo structs.EndpointInfo
	conn.ReadJSON(&endpointInfo)

	endPoint := models.MakeEndpoint(endpointInfo)
	hubconn := newConnection(conn, endPoint)
	defer func() {
		hub.deleteConnection(hubconn.Id)
		conn.Close()
	}()
	fmt.Printf("Got a new endpoint:, %v", req.Host)
}

func (hub *Hub) addConnection(conn *_Connection) {
	hub.mx.Lock()
	hub.endpoints[conn.Id] = conn
	hub.mx.Unlock()
}

func (hub *Hub) deleteConnection(id structs.UUID) {
	hub.mx.Lock()
	delete(hub.endpoints, id)
	hub.mx.Unlock()
}
