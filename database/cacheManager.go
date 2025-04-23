package database

import (
	"errors"
	"os"

	// "fmt"
	// "log"
	"encoding/json"
)

type Cache struct {
	UserCount       int `json:"user_count"`
	ProblemCount    int `json:"problem_count"`
	SubmissionCount int `json:"submission_count"`
}

var cacheFileName string = "cache.json"

func readCache(path string) (Cache, error) {
	var cache Cache

	file, err := os.Open(path)
	if err != nil {
		return cache, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return cache, err
	}

	data := make([]byte, info.Size())
	bytesRead, err := file.Read(data)
	if err != nil {
		return cache, err
	}
	if int64(bytesRead) != info.Size() {
		return cache, errors.New("could not read entire file")
	}

	err = json.Unmarshal(data, &cache)
	return cache, err
}

func writeCache(path string, cache Cache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
