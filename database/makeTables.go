package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Is_admin           bool   `json:"is_admin"`
	Attempted_problems int    `json:"attempted_problems"`
	Solved_problems    int    `json:"solved_problems"`
}

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

func connect() (*sql.DB, error) {
	var err error

	username := "Mahdi"           // your MySQL username
	password := "Mahdi0441265367" // your MySQL password
	hostname := "127.0.0.1:3306"
	dbname := "codehub" // your MySQL database name

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbname)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	makeUsers := `
    create table if not exists users (
		id integer not null primary key AUTO_INCREMENT, 
		username varchar(50) not null unique, 
		password varchar(255) not null,
		is_admin BOOLEAN DEFAULT false,
		attempted_problems INT DEFAULT 0, 
		solved_problems INT DEFAULT 0
	);
    `
	makeProblems := `
	CREATE TABLE IF NOT EXISTS problems (
		id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, 
		user_id INT NOT NULL, 
		title VARCHAR(50) NOT NULL, 
		description_path VARCHAR(255) NOT NULL,
		input_path VARCHAR(255) NOT NULL,
		output_path VARCHAR(255) NOT NULL, 
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		is_published BOOLEAN DEFAULT false, 
		time_limit_ms INT,
		memory_limit_mb INT 
	);
	`

	makeSubmissions := `
	CREATE TABLE IF NOT EXISTS submissions (
		id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, 
		user_id INT NOT NULL,
		problem_id INT NOT NULL, 
		code_path VARCHAR(255) NOT NULL,
		state TINYINT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		runtime_ms INT DEFAULT 0,
		memory_used INT DEFAULT 0,
		error_message TEXT DEFAULT ''
	);
	`

	/*

		states in submission
		0 -> pending
		1 -> OK
		2 -> compile error
		3 -> wrong answer
		4 -> memory limit
		6 -> time limit
		7 -> runtime error

	*/
	_, err = db.Exec(makeUsers)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(makeProblems)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(makeSubmissions)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	_, err := connect()
	if err != nil {
		panic(err)
	}
}
