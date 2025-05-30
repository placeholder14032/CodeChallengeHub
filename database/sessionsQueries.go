package database

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
    "database/sql"

    "github.com/placeHolder143032/CodeChallengeHub/models"

	// "time"
	_ "github.com/go-sql-driver/mysql"
    "log"
)

func GenerateSessionID() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(userID int) (string, error) {
    sessionID, err := GenerateSessionID()
    if err != nil {
        fmt.Println("Error generating session ID:", err)
        return "", err
    }

    // Create sessions table if it doesn't exist
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            session_id VARCHAR(64) PRIMARY KEY,
            user_id INTEGER NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            expires_at TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
    `)
    if err != nil {
        fmt.Println("Error creating sessions table:", err)
        return "", err
    }

    // Insert new session
    _, err = db.Exec(`
        INSERT INTO sessions (session_id, user_id, expires_at) 
        VALUES (?, ?, DATE_ADD(NOW(), INTERVAL 24 HOUR))
    `, sessionID, userID)
    if err != nil {
        fmt.Println("Error inserting session into database:", err)
        return "", err
    }

    return sessionID, nil
}

func ValidateSession(sessionID string) (int, bool, error) {
    var id int
    err := db.QueryRow(`
        SELECT user_id 
        FROM sessions 
        WHERE session_id = ? AND expires_at > NOW()
    `, sessionID).Scan(&id)

    if err == sql.ErrNoRows {
        log.Printf("Session %s not found or expired", sessionID)
        return 0, false, nil
    }
    if err != nil {
        log.Printf("Database error validating session %s: %v", sessionID, err)
        return 0, false, err
    }

    // Refresh session expiration
    _, err = db.Exec(`
        UPDATE sessions 
        SET expires_at = DATE_ADD(NOW(), INTERVAL 24 HOUR)
        WHERE session_id = ?
    `, sessionID)
    if err != nil {
        log.Printf("Error refreshing session %s: %v", sessionID, err)
    }

    log.Printf("Session %s validated, user_id=%d", sessionID, id)
    return id, true, nil
}


func GetCurrentUser(sessionID string) (*models.User, error) {
    // First, get user_id from sessions table
    var userID int
    err := db.QueryRow(`
        SELECT user_id FROM sessions 
        WHERE session_id = ? AND expires_at > NOW()
    `, sessionID).Scan(&userID)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("invalid or expired session")
        }
        return nil, err
    }

    // Then get user details
    var user models.User
    err = db.QueryRow(`
        SELECT id, username, attempted_problems, solved_problems, is_admin 
        FROM users 
        WHERE id = ?
    `, userID).Scan(&user.ID, &user.Username, &user.AttemptedProblems, &user.SolvedProblems, &user.IsAdmin)
    
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func GetCurrentUserID(sessionID string) (*models.User, error) {
    // First, get user_id from sessions table
    var userID int
    err := db.QueryRow(`
        SELECT user_id FROM sessions 
        WHERE session_id = ? AND expires_at > NOW()
    `, sessionID).Scan(&userID)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("invalid or expired session")
        }
        return nil, err
    }

    // Then get user details
    var user models.User
    err = db.QueryRow(`
        SELECT id, username, attempted_problems, solved_problems, is_admin 
        FROM users 
        WHERE id = ?
    `, userID).Scan(&user.ID, &user.Username, &user.AttemptedProblems, &user.SolvedProblems, &user.IsAdmin)
    
    if err != nil {
        return nil, err
    }

    return &user, nil
}


