package models

import (
	"database/sql"
	"time"
)

type Submission struct {
	ID            int            `json:"id"`
	UserId        int            `json:"user_id"`
	ProblemId     int            `json:"problem_id"`
	ProblemTitle  string         `json:"problem_title"`
	CodePath      string         `json:"code_path"`
	State        int       `json:"state"`      // Numeric state from database
    Status       string    `json:"status"`     // String representation of state
	CreatedAt     time.Time      `json:"created_at"`
	Runtime_ms    int            `json:"runtime_ms"`
	Memory_used   int            `json:"memory_used"`
	Error_message sql.NullString `json:"error_message"`
	JudgeID 	int64 
}