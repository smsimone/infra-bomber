package queue

type Task struct {
	Input    map[string]string
	FlowFile string
}

func NewTask(flowFile string, input map[string]string) *Task {
	return &Task{
		FlowFile: flowFile,
		Input:    input,
	}
}
