package models

import(
	"time"
)

type Problem struct {
	ID               int       `json:"id"`
	User_id          int       `json:"user_id"`
	Title            string    `json:"title"`
	Description_path string    `json:"description_path"`
	Input_path       string    `json:"input_path"`
	Output_path      string    `json:"output_path"`
	Created_at       time.Time `json:"created_at"`
	Is_Published     bool      `json:"is_published"`
	Time_limit_ms    int       `json:"time_limit_ms"`
	Memory_limit_mb  int       `json:"memory_limit_mb"`
}