package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

func SignUpUser(user User) error {
	var exists bool
	query := "SELECT 1 FROM users WHERE username = ? LIMIT 1"
	err := db.QueryRow(query, user.Username).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if exists {
		return errors.New("user already exists")
	}

	insertQuery := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = db.Exec(insertQuery, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func SignInUser(user User) (int, error) {
	var id int
	var password string
	query := "SELECT id, password FROM users WHERE username = ? LIMIT 1"
	err := db.QueryRow(query, user.Username).Scan(&id, &password)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, errors.New("user does not exist")
	}

	if password != user.Password {
		return 0, errors.New("wring password")
	}

	return id, nil
}

func GetUserInfo(id int) (User, error) {
	var user User
	query := `
	SELECT username, password, is_admin, attempted_problems, solved_problems FROM users 
	WHERE id = ? LIMIT 1
	`
	err := db.QueryRow(query, id).Scan(&user.Username, &user.Password, &user.Is_admin, &user.Attempted_problems, &user.Solved_problems)

	if err != nil && err != sql.ErrNoRows {
		return User{}, err
	}

	if err == sql.ErrNoRows {
		return User{}, errors.New("user does not exist")
	}

	return user, nil
}

func GetUserRole(id int) (bool, error) {
	var is_admin bool
	query := `
	SELECT is_admin FROM users 
	WHERE id = ? LIMIT 1
	`
	err := db.QueryRow(query, id).Scan(&is_admin)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, errors.New("user does not exist")
	}

	return is_admin, nil
}

func ChangeUserRole(id int) error {
	query := `
	UPDATE users SET is_admin = NOT is_admin 
	WHERE id = ? LIMIT 1
	`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func ChangeUserPassword(id int, password string) error {
	query := `
	UPDATE users SET password = ?
	WHERE id = ? LIMIT 1
	`
	result, err := db.Exec(query, password, id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user does not exist")
	}

	if err != nil {
		return err
	}

	return nil
}
