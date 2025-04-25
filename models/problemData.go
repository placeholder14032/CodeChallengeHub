package models

import(
	// "time"
)

type ProblemData struct {
	ID               int       `json:"id"`
	UserID          int       `json:"user_id"`
	Title            string    `json:"title"`
	Explanation string    `json:"description"`
	Input      string    `json:"input"`
	Output      string    `json:"output"`
	// CreatedTime       time.Time `json:"created_time"`
	IsPublished     bool      `json:"is_published"`
	TimeLimit	   int       `json:"time_limit"`
	MemoryLimit  	int       `json:"memory_limit"`
}