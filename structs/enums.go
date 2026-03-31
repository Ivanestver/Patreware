package structs

type EndpointSecurityState int

const (
	SecurityStateInfected EndpointSecurityState = iota
	SecurityStateClean
	SecurityStateSuspicious
)

type EndpointConnectionState int

const (
	ConnectionStateConnected EndpointConnectionState = iota
	ConnectionStateNotConnected
)
