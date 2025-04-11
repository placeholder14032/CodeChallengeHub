package database

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	// "time"
	_ "github.com/go-sql-driver/mysql"
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

func ValidateSession(sessionID string) (bool, error) {
    var exists bool
    err := db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM sessions 
            WHERE session_id = ? AND expires_at > NOW()
        )
    `, sessionID).Scan(&exists)

    if err != nil {
        return false, err
    }

    return exists, nil
}