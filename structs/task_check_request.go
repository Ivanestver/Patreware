package structs

import "time"

type TaskCheckRequest struct {
	TaskBase
	Time time.Time `json:"datetime"`
}

func MakeTaskCheckRequest(taskname string) TaskCheckRequest {
	return TaskCheckRequest{
		TaskBase: MakeTaskBase(taskname),
		Time:     time.Now(),
	}
}
