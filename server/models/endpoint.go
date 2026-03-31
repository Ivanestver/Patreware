/*
Package models is blah-blah-blah
*/
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
		connectionState: structs.ConnectionStateConnected,
	}
}

func (endpoint *Endpoint) GetID() structs.UUID {
	return endpoint.id
}

func (endpoint *Endpoint) IsConnected() bool {
	return endpoint.connectionState == structs.ConnectionStateConnected
}
