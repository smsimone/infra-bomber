package queue

import (
	"fmt"
	"it.toduba/bomber/flow"
	"log"
)

type Queue struct {
	Tasks            []Task
	currentlyRunning int
	MaxExecutions    int
	onRunningChange  chan int
}

type AddProp func(q *Queue)

// Initialize the Queue with some custom props
func (q *Queue) Initialize(props ...AddProp) {
	for _, p := range props {
		p(q)
	}

	q.onRunningChange = make(chan int, q.MaxExecutions)
}

// AddTask appends a new task t into the queue
func (q *Queue) AddTask(t *Task) {
	q.Tasks = append(q.Tasks, *t)
}

// popTask Returns the first task in queue. Returns nil if empty
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

// runNextTask pops the first item in the queue and runs it
func (q *Queue) runNextTask() error {
	canRun := false

	for !canRun {
		select {
		case <-q.onRunningChange:
			if q.currentlyRunning < q.MaxExecutions {
				q.currentlyRunning += 1
				q.onRunningChange <- q.currentlyRunning
				canRun = true
				break
			}
		}
	}

	fmt.Printf("Currently running %v jobs\n", q.currentlyRunning)

	t := q.popTask()
	if t == nil {
		return fmt.Errorf("queue has no tasks in it")
	}

	go func(task *Task) {

		defer func() {
			log.Printf("Releasing task")
			q.currentlyRunning -= 1
			q.onRunningChange <- q.currentlyRunning
		}()

		f, err := flow.ParseFromYaml(task.FlowFile)
		if err != nil {
			log.Printf("Failed to parse flow: %v\n", err.Error())
		}

		if err := f.Execute((*task).Input); err != nil {
			fmt.Printf("Failed to execute task: %v\n", err.Error())
		}
	}(t)

	return nil
}

// Start Starts to execute all tasks
func (q *Queue) Start() {
	q.onRunningChange <- -1
	for len(q.Tasks) > 0 {
		if err := q.runNextTask(); err != nil {
			log.Printf("Failed to execute next task: %v\n", err.Error())
		}
	}
}

// Wait Let the process wait to empty the queue
func (q *Queue) Wait() {
	for {
		select {
		case <-q.onRunningChange:
			fmt.Printf("----- Currently there are %v processes\n", q.currentlyRunning)
			if q.currentlyRunning == 0 {
				fmt.Printf("------ exiting\n")
				return
			}
		}
	}
}
