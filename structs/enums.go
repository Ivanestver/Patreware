package structs

type EndpointSecurityState int

const (
	SECURITY_STATE_INFECTED EndpointSecurityState = iota
	SECURITY_STATE_CLEAN
	SECURITY_STATE_SUSPICIOUS
)

type EndpointConnectionState int

const (
	CONNECTION_STATE_CONNECTED EndpointConnectionState = iota
	CONNECTION_STATE_NOT_CONNECTED
)
