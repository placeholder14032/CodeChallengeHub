package models

import(
	"time"
)

type Submission struct {
	ID            int       `json:"id"`
	User_id       int       `json:"user_id"`
	Problem_id    int       `json:"problem_id"`
	Code_path     string    `json:"code_path"`
	State         int8      `json:"state"`
	Created_at    time.Time `json:"created_at"`
	Runtime_ms    int       `json:"runtime_ms"`
	Memory_used   int       `json:"memory_used"`
	Error_message string    `json:"error_message"`
}