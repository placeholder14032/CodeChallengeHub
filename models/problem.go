package models

import(
	"time"
)

type Problem struct {
	ID               int       `json:"id"`
	UserID          int       `json:"user_id"`
	Title            string    `json:"title"`
	DescriptionPath string    `json:"description_path"`
	InputPath       string    `json:"input_path"`
	OutputPath      string    `json:"output_path"`
	CreatedTime       time.Time `json:"created_time"`
	IsPublished     bool      `json:"is_published"`
	TimeLimit	   int       `json:"time_limit"`
	MemoryLimit  	int       `json:"memory_limit"`
}