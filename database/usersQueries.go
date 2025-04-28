package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/placeHolder143032/CodeChallengeHub/models"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

func SignUpUser(user models.User) error {
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

func SignInUser(user models.User) (int, string, error) {
    var id int
    var password string
    query := "SELECT id, password FROM users WHERE username = ? LIMIT 1"
    err := db.QueryRow(query, user.Username).Scan(&id, &password)

    if err != nil && err != sql.ErrNoRows {
        return 0, "", err
    }

    if err == sql.ErrNoRows {
        return 0, "", errors.New("user does not exist")
    }

    if !VerifyPassword(user.Password, password) {
        return 0, "", errors.New("wrong password")
    }

    // Create session
    sessionID, err := CreateSession(id)
    if err != nil {
        return 0, "", errors.New("Failed to create session")
    }

    return id, sessionID, nil
}

func VerifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func GetUserInfo(id int) (models.User, error) {
    var user models.User
    query := `
    SELECT id, username, password, is_admin, 
           COALESCE(attempted_problems, 0) as attempted_problems, 
           COALESCE(solved_problems, 0) as solved_problems 
    FROM users 
    WHERE id = ? 
    LIMIT 1
    `
    err := db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Password,
        &user.IsAdmin,
        &user.AttemptedProblems,
        &user.SolvedProblems,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("User not found: id=%d", id)
            return models.User{}, errors.New("user does not exist")
        }
        log.Printf("Database error in GetUserInfo: id=%d, error=%v", id, err)
        return models.User{}, fmt.Errorf("database error: %w", err)
    }

    log.Printf("Successfully retrieved user: id=%d, username=%s", user.ID, user.Username)
    return user, nil
}

func GetUserRole(id int) (int, error) {
	var is_admin int
	query := `
	SELECT is_admin FROM users 
	WHERE id = ? LIMIT 1
	`
	err := db.QueryRow(query, id).Scan(&is_admin)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, errors.New("user does not exist")
	}

	return is_admin, nil
}


func MakeAdmin(id int) error {
	query := `
	UPDATE users
	SET is_admin = 1
	WHERE id = ? LIMIT 1
	`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

func RemoveAdmin(id int) error {
	query := `
	UPDATE users
	SET is_admin = 0
	WHERE id = ? LIMIT 1
	`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

func ChangeUserPassword(id int, password string) error {
	query := `
	UPDATE users SET password = ?
	WHERE id = ? LIMIT 1
	`
	result, err := db.Exec(query, password, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

func UpdateUserProblemStats() error {
	query := `
		UPDATE users u
		SET 
			attempted_problems = stats.attempted,
			solved_problems = stats.solved
		FROM (
			SELECT 
				user_id,
				COUNT(DISTINCT problem_id) AS attempted,
				COUNT(DISTINCT CASE WHEN state = 1 THEN problem_id END) AS solved
			FROM submissions
			GROUP BY user_id
		) stats
		WHERE u.id = stats.user_id;
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}


func CreateAdminUser(username, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %v", err)
    }

	fmt.Println("print in create admin user befor query")

    query := `
        INSERT INTO users (username, password, is_admin)
        VALUES (?, ?, 1)
    `

    _, err = db.Exec(query, username, hashedPassword)
    if err != nil {
        return fmt.Errorf("failed to create admin user: %v", err)
    }

	fmt.Println("print in create admin after befor query yayyyy")


    return nil
}


func GetUserIDFromSession(sessionID string) (int, error) {
    var userID int
    query := `
        SELECT user_id FROM sessions 
        WHERE session_id = ? AND expires_at > NOW()
        LIMIT 1
    `
    err := db.QueryRow(query, sessionID).Scan(&userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, fmt.Errorf("session not found or expired")
        }
        return 0, fmt.Errorf("database error: %v", err)
    }

    log.Printf("Found user_id=%d for session=%s", userID, sessionID)
    return userID, nil
}