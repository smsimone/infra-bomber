package queue

import (
	"log"

	"it.toduba/bomber/flow"
)

type Task struct {
	BaseTask
	Input    *map[string]string
	FlowFile string
}

func NewTask(flowFile string, input *map[string]string) *Task {
	return &Task{
		FlowFile: flowFile,
		Input:    input,
	}
}

func (t *Task) Execute() error {
	f, err := flow.ParseFromYaml(t.FlowFile)
	if err != nil {
		log.Printf("Failed to parse flow: %v\n", err.Error())
	}

	return f.Execute((*t).Input)
}
