package models

type User struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Is_admin           int    `json:"is_admin"`
	Attempted_problems int    `json:"attempted_problems"`
	Solved_problems    int    `json:"solved_problems"`
}