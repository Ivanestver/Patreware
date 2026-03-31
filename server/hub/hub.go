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

var _HubInstance *Hub

func InitHub() {
	_HubInstance = &Hub{
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
	_HubInstance.setupHandlers()
}

func (hub *Hub) setupHandlers() {
	http.HandleFunc("/ws", hub.serveConnection)
}

func (hub *Hub) serveConnection(writer http.ResponseWriter, req *http.Request) {
	conn, err := hub.upgrader.Upgrade(writer, req, nil)
	if err != nil {
		panic("Not implemented")
	}
	var endpointInfo structs.EndpointInfo
	conn.ReadJSON(&endpointInfo)

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

func GetConnectionAssisiatedWithEndpoint(endpointId structs.UUID) (*_Connection, error) {
	for _, conn := range _HubInstance.endpoints {
		if conn.Endpoint.GetID().Equals(endpointId) {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("no connection assiciated with the endpoint: ", endpointId)
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
