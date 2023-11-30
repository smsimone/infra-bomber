package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"it.toduba/bomber/queue"
	"it.toduba/bomber/utils"
)

func main() {
	f := parseFlags()

	q := queue.Queue{}
	q.Initialize(
		func(q *queue.Queue) {
			q.MaxExecutions = f.Jobs
		},
		func(q *queue.Queue) {
			for _, t := range GetTasks(*f) {
				q.AddTask(&t)
			}
		},
	)

	log.Printf("Should run %v iterations\n", len(q.Tasks))

	q.Start()
	q.Wait()
}

func GetTasks(f utils.ProgramArgs) []queue.BaseTask {
	var vars []map[string]string
	if f.Variables != nil {
		vars = utils.ReadInputCsv(*f.Variables)
	}

	var files []string

	if stat, _ := os.Stat(f.Flow); stat.IsDir() {
		if entries, err := os.ReadDir(f.Flow); err != nil {
			panic(fmt.Sprintf("Failed to read flow dir: %v", err.Error()))
		} else {
			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml") {
					files = append(files, path.Join(f.Flow, entry.Name()))
				}
			}
		}
	} else {
		files = append(files, f.Flow)
	}

	var tasks []queue.BaseTask
	if len(vars) == 0 {
		for _, f := range files {
			tasks = append(tasks, queue.NewTask(f, nil))
		}
	} else {
		count := 0
		for _, group := range vars {
			if *f.Limit != -1 && count >= *f.Limit {
				break
			}
			for _, f := range files {
				tasks = append(tasks, queue.NewTask(f, &group))
			}
			count += 1
		}
	}

	return tasks
}

func parseFlags() *utils.ProgramArgs {
	jobs := flag.Int("jobs", 10, "Numero di job concorrenti da eseguire")
	flow := flag.String("flow", "", "Percorso al flusso da eseguire. Se random è true, deve puntare ad una cartella")
	variables := flag.String("variables", "", "Percorso al file csv contenente le variabili per sostituire i placeholder nel flusso")
	limit := flag.Int("limit", -1, "Specifica il limite di iterazioni da eseguire (solo se variables è definito)")
	flag.Parse()

	if len(*flow) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*flow); os.IsNotExist(err) {
		panic("The flow path does not exists")
	}

	if *jobs <= 0 {
		panic("Invalid job number")
	}

	if len(*variables) != 0 {
		if stat, err := os.Stat(*variables); os.IsNotExist(err) || stat.IsDir() {
			panic("Variables path points to a non existent file or to a directory")
		}
		if !strings.HasSuffix(*variables, ".csv") {
			panic("Invalid variables file")
		}
	}

	return &utils.ProgramArgs{
		Jobs:      *jobs,
		Flow:      *flow,
		Variables: variables,
		Limit:     limit,
	}
}
