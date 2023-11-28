package queue

import (
	"fmt"
	"log"
)

type Queue struct {
	onRunningChange  chan int
	Tasks            []BaseTask
	currentlyRunning int
	MaxExecutions    int
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
func (q *Queue) AddTask(t *BaseTask) {
	q.Tasks = append(q.Tasks, *t)
}

// popTask Returns the first task in queue. Returns nil if empty
func (q *Queue) popTask() *BaseTask {
	if len(q.Tasks) == 0 {
		return nil
	}

	if len(q.Tasks) == 1 {
		t := q.Tasks[0]
		q.Tasks = []BaseTask{}
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
		<-q.onRunningChange
		if q.currentlyRunning < q.MaxExecutions {
			q.currentlyRunning += 1
			q.onRunningChange <- q.currentlyRunning
			canRun = true
		}
	}

	fmt.Printf("Currently running %v jobs\n", q.currentlyRunning)

	t := q.popTask()
	if t == nil {
		return fmt.Errorf("queue has no tasks in it")
	}

	go func(task *BaseTask) {
		defer func() {
			q.currentlyRunning -= 1
			q.onRunningChange <- q.currentlyRunning
		}()

		if err := (*task).Execute(); err != nil {
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
		<-q.onRunningChange
		fmt.Printf("----- Currently there are %v processes\n", q.currentlyRunning)
		if q.currentlyRunning == 0 {
			fmt.Printf("------ exiting\n")
			return
		}
	}
}
