package services

import (
	"encoding/csv"
	"os"
	"strconv"
)

type ReaderResponse []int

func readFile(path string) (*ReaderResponse, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	results, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var response ReaderResponse

	// loop rows
	for _, r := range results {
		// first column
		res, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, err
		}
		response = append(response, res)
	}

	return &response, nil
}
