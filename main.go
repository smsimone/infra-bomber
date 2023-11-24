package main

import (
	"fmt"
	"it.toduba/bomber/queue"
)

func main() {
	vars := ReadInputCsv("resources/variables.csv")

	q := queue.Queue{}
	q.Initialize(
		func(q *queue.Queue) {
			q.MaxExecutions = 20
		},
		func(q *queue.Queue) {
			for _, group := range vars {
				q.AddTask(queue.NewTask("resources/sample_flow.yaml", group))
			}
		},
	)

	fmt.Printf("Should run %v iterations", len(q.Tasks))

	q.Start()
}
