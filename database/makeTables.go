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
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
        username VARCHAR(50) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        isAdmin BOOLEAN DEFAULT 0,
        attemptedProblems INT DEFAULT 0,
        solvedProblems INT DEFAULT 0
    );
    `

	makeProblems := `
    CREATE TABLE IF NOT EXISTS problems (
        id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
        user_id INT NOT NULL,
        title VARCHAR(50) NOT NULL,
        difficulty VARCHAR(20) NOT NULL,
        description_path VARCHAR(255) NOT NULL,
        input_path VARCHAR(255) NOT NULL,
        output_path VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_published BOOLEAN DEFAULT FALSE,
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

	makeSessions := `
    CREATE TABLE IF NOT EXISTS sessions (
        session_id VARCHAR(64) PRIMARY KEY,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );
    `

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
	_, err = db.Exec(makeSessions)
	if err != nil {
		return nil, err
	}
	return db, nil
}

