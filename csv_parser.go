package main

import (
	"encoding/csv"
	"it.toduba/bomber/utils"
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
	utils.HandleErrFail(err)

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	reader := csv.NewReader(f)
	data, err := reader.ReadAll()
	utils.HandleErrFail(err)
	return parseData(data)
}
