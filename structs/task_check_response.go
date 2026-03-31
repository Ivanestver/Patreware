package structs

import "net"

type TaskCheckResponse struct {
	TaskBase
	IP            net.IPNet             `json:"ip"`
	SecurityState EndpointSecurityState `json:"security_state"`
}

func MakeTaskCheckResponse(taskname string, ip net.IPNet) TaskCheckResponse {
	return TaskCheckResponse{
		TaskBase:      MakeTaskBase(taskname),
		IP:            ip,
		SecurityState: SecurityStateClean,
	}
}
