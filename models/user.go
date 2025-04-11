package models

type User struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	IsAdmin           int    `json:"is_admin"`
	AttemptedProblems int    `json:"attempted_problems"`
	SolvedProblems    int    `json:"solved_problems"`
	SuccessRate	  float64 `json:"success_rate"`
}