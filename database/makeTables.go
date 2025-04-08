package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func connect() (*sql.DB, error) {
	var err error

	username := "Mahdi"
	password := "Mahdi0441265367"
	hostname := "127.0.0.1:3306"
	dbname := "codehub"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbname)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	makeUsers := `
    create table if not exists users (
		id integer not null primary key AUTO_INCREMENT, 
		username varchar(50) not null unique, 
		password varchar(255) not null
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
		runtime_ms INT,
		memory_used INT,
		error_message TEXT
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
