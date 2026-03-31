package structs

type TaskBase struct {
	Name   string `json:"name"`
	TaskID UUID
}

func MakeTaskBase(taskname string) TaskBase {
	return TaskBase{
		Name:   taskname,
		TaskID: GenerateUUID(),
	}
}
