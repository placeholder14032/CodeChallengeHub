package database

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

// execute before everything because if not might choose deleted users/problems as the random user/problem

func loadTest() error {
	err := addUser(1000)
	if err != nil {
		return err
	}

	err = addProblem(100)
	if err != nil {
		return err
	}

	err = addSubmission(250)
	if err != nil {
		return err
	}

	return nil
}

func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"
	lenCharset := len(charset)
	password := make([]byte, length)

	for i := range password {
		n := rand.IntN(lenCharset)
		password[i] = charset[n]
	}

	return string(password), nil
}

func addUser(count int) error {
	cache, _ := readCache(cacheFileName)
	user_count := cache.UserCount
	for i := 0; i < count; i++ {
		username := fmt.Sprintf("username%d", user_count+1)
		length := rand.IntN(9) + 8 // generating length between in [8,16]
		password, err := generateRandomPassword(length)
		if err != nil {
			return errors.New("error in creating users: password")
		}
		SignInUser(User{Username: username, Password: password})
	}
	return nil
}

func addProblem(count int) error {
	cache, _ := readCache(cacheFileName)
	user_count, problem_count := cache.UserCount, cache.ProblemCount
	for i := 0; i < count; i++ {
		user_idRandom := rand.IntN(user_count) + 1
		// title, created_at, time_limit, memory_limit
		// creating files for input, output & description
		// creates file with (problem_count+1). this gonna be the id of the problem

		// File path
		descPath := fmt.Sprintf("problems/desc/%d.txt", problem_count+1)
		inputPath := fmt.Sprintf("problems/input/%d.txt", problem_count+1)
		outputPath := fmt.Sprintf("problems/output/%d.txt", problem_count+1)

		// File content
		inputContent := fmt.Sprintf("input for problem %d", problem_count+1)
		outputContent := fmt.Sprintf("output for problem %d", problem_count+1)
		descriptionContent := fmt.Sprintf("This is problem number %d about a GO project", problem_count+1)

		// Write files
		if err := createAndWrite(descPath, descriptionContent); err != nil {
			return err
		}
		if err := createAndWrite(inputPath, inputContent); err != nil {
			return err
		}
		if err := createAndWrite(outputPath, outputContent); err != nil {
			return err
		}
		title, created_at, time_limit_ms, memory_limit_mb := fmt.Sprintf("title of problem %d", problem_count+1), time.Now(), 1000, 256

		if err := AddProblem(user_idRandom,
			Problem{Title: title,
				Description_path: descPath,
				Input_path:       inputPath,
				Output_path:      outputPath,
				Created_at:       created_at,
				Time_limit_ms:    time_limit_ms,
				Memory_limit_mb:  memory_limit_mb}); err != nil {
			return err
		}

	}
	return nil
}

func createAndWrite(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	return err
}

func addSubmission(count int) error {
	cache, _ := readCache(cacheFileName)
	user_count, problem_count, submission_count := cache.UserCount, cache.ProblemCount, cache.SubmissionCount
	for i := 0; i < count; i++ {
		user_idRandom, problem_idRandom := rand.IntN(user_count)+1, rand.IntN(problem_count)+1

		// file path
		codePath := fmt.Sprintf("problems/desc/%d.txt", submission_count+1)

		// file content
		code := fmt.Sprintf("package main\n\nfunc main() {\n\tprintln(\"Hello from submission %d\")\n}\n", submission_count+1)

		// write file
		if err := createAndWrite(codePath, code); err != nil {
			return err
		}
		created_at := time.Now()

		if err := SubmitCode(Submission{User_id: user_idRandom, Problem_id: problem_idRandom, Code_path: codePath, Created_at: created_at}); err != nil {
			return err
		}
	}
	return nil
}
