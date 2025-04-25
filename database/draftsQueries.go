package database

import (
	"fmt"
)

func ToggleProblemDraftStatus(problemID int) error {
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