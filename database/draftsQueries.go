package database

import (
	"database/sql"
	"fmt"
)

func ToggleProblemDraftStatus(problemID int) error {
	// Check if the problem exists
	query1 := `
		SELECT id 
		FROM problems 
		WHERE id = ?
	`

	var id int
	err := db.QueryRow(query1, problemID).Scan(&id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("problem not found")
	}
	if err != nil {
		return err
	}

	// Toggle is_published: TRUE -> FALSE, FALSE -> TRUE
	query := `
		UPDATE problems 
		SET is_published = NOT is_published
		WHERE id = ?
	`
	
	result, err := db.Exec(query, problemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}