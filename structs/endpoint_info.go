package structs

type EndpointInfo struct {
	Name          string                `json:"name"`
	SecurityState EndpointSecurityState `json:"state"`
}
