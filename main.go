package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"it.toduba/bomber/flow"
)

func main() {
	f, err := flow.ParseFromYaml("resources/sample_flow.yaml")
	if err != nil {
		log.Fatalf("Failed to parse flow: %v", err.Error())
	}
	log.Printf("Read schema: %v", f.Name)

	content := readInputCsv("resources/variables.csv")
	vars := parseData(content)

	total := len(vars)
	for idx, group := range vars {
		f.Execute(&group)
		fmt.Printf("Done group %v/%v\n", idx, total)
		break
	}
}

func parseData(csvData [][]string) []map[string]string {
	header := csvData[0]

	vars := make([]map[string]string, 0)
	for idx, line := range csvData {
		if idx == 0 {
			continue
		}
		group := make(map[string]string)
		for hIdx, h := range header {
			group[h] = line[hIdx]
		}
		vars = append(vars, group)
	}

	return vars
}

func readInputCsv(csvFile string) [][]string {
	f, err := os.Open(csvFile)
	handleErr(err)

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	reader := csv.NewReader(f)
	data, err := reader.ReadAll()
	handleErr(err)
	return data
}

func handleErr(err error) {
	if err != nil {
		log.Fatalf("Got error: %v", err)
	}
}
