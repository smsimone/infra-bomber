package utils

import (
	"encoding/csv"
	"log"
	"os"
)

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

func ReadInputCsv(csvFile string) []map[string]string {
	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open variables file: %v", err.Error())
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	reader := csv.NewReader(f)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read variables file: %v", err.Error())
	}
	return parseData(data)
}
