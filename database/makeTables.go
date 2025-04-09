package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect() (*sql.DB, error) {
	var err error

	username := "amabilis"           // your MySQL username
	password := "amabilisfi20050921" // your MySQL password
	hostname := "127.0.0.1:3306"
	dbname := "codeChallemgeHub" // your MySQL database name

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
		is_admin BOOLEAN DEFAULT 0,
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
	/*

		states in is_admin
		0 -> simple user
		1 -> user with admin access
		2 -> admin

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