package models

import (
	"patrware/structs"
)

type Endpoint struct {
	structs.EndpointInfo
	id              structs.UUID
	connectionState structs.EndpointConnectionState
}

func MakeEndpoint(endpointInfo structs.EndpointInfo) Endpoint {
	return Endpoint{
		EndpointInfo:    endpointInfo,
		id:              structs.GenerateUUID(),
		connectionState: structs.CONNECTION_STATE_CONNECTED,
	}
}

func (endpoint *Endpoint) GetID() structs.UUID {
	return endpoint.id
}

func (endpoint *Endpoint) IsConnected() bool {
	return endpoint.connectionState == structs.CONNECTION_STATE_CONNECTED
}
