package queue

import (
	"fmt"
	"it.toduba/bomber/flow"
	"log"
	"sync"
)

type Queue struct {
	Tasks          []Task
	mutex          sync.Mutex
	currentRunning int
	MaxExecutions  int
}

type AddProp func(q *Queue)

func (q *Queue) Initialize(props ...AddProp) {
	q.mutex = sync.Mutex{}
	for _, p := range props {
		p(q)
	}
}

func (q *Queue) AddTask(t *Task) {
	q.Tasks = append(q.Tasks, *t)
}

func (q *Queue) releaseTask() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.currentRunning -= 1
}

func (q *Queue) popTask() *Task {
	if len(q.Tasks) == 0 {
		return nil
	}

	if len(q.Tasks) == 1 {
		t := q.Tasks[0]
		q.Tasks = []Task{}
		return &t
	}

	t := q.Tasks[0]
	q.Tasks = q.Tasks[1:]

	return &t
}

func (q *Queue) RunNextTask() error {
	canRun := false

	for !canRun {
		q.mutex.Lock()
		func() {
			defer q.mutex.Unlock()
			if q.currentRunning < q.MaxExecutions {
				q.currentRunning += 1
				canRun = true
			}
		}()
	}

	fmt.Printf("Currently running %v jobs", q.currentRunning)

	t := q.popTask()
	if t == nil {
		return fmt.Errorf("queue has no tasks in it")
	}
	go func(task *Task) {
		defer q.releaseTask()

		f, err := flow.ParseFromYaml("resources/sample_flow.yaml")
		if err != nil {
			log.Fatalf("Failed to parse flow: %v", err.Error())
		}

		f.Execute(&(*task).Input)
	}(t)
	return nil
}

// Start Starts to execute all tasks
func (q *Queue) Start() {
	for len(q.Tasks) > 0 {
		if err := q.RunNextTask(); err != nil {
			log.Fatalf("Failed to execute next task: %v", err.Error())
		}
	}
}
