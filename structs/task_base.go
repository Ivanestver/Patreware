package structs

type TaskBase struct {
	Name   string `json:"name"`
	TaskId UUID
}

func MakeTaskBase(taskname string) TaskBase {
	return TaskBase{
		Name:   taskname,
		TaskId: GenerateUUID(),
	}
}
