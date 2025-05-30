package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/placeHolder143032/CodeChallengeHub/models"

	_ "github.com/go-sql-driver/mysql"
)

// this one is for admins we womt check if its published or not anymore
func GetProblemsPageAdmin(m, n int) ([]models.Problem, error) {
    var problems []models.Problem
    query := `
        SELECT id, title, created_at, is_published
        FROM problems
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
    offset := (m - 1) * n

    rows, err := db.Query(query, n, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var p models.Problem
        if err := rows.Scan(&p.ID, &p.Title, &p.CreatedTime, &p.IsPublished); err != nil {
            return nil, err
        }
        problems = append(problems, p)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return problems, nil
}

// this one is for users since it does check if its published or not
func GetProblemsPageUser(m, n int) ([]models.Problem, error) {
    var problems []models.Problem
    query := `
        SELECT 
            id, 
            user_id,
            title, 
            created_at,
            is_published
        FROM problems
        WHERE is_published = true
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
    offset := (m - 1) * n

    rows, err := db.Query(query, n, offset)
    if err != nil {
        log.Printf("Query error: %v", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var p models.Problem
        if err := rows.Scan(
            &p.ID,
            &p.UserID,
            &p.Title,
            &p.CreatedTime,
            &p.IsPublished,
        ); err != nil {
            log.Printf("Scan error: %v", err)
            return nil, err
        }
        problems = append(problems, p)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Rows error: %v", err)
        return nil, err
    }

    return problems, nil
}

// created at should be from the other side for less inconsistency
func AddProblem(userID int, problem models.Problem) error {
    query := `
        INSERT INTO problems (
            user_id,
            title,
            description_path,
            input_path,
            output_path,
            created_at,
            is_published,
            time_limit_ms,
            memory_limit_mb
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    _, err := db.Exec(query,
        problem.UserID,
        problem.Title,
        problem.DescriptionPath,
        problem.InputPath,
        problem.OutputPath,
        problem.CreatedTime,
        problem.IsPublished,
        problem.TimeLimit,
        problem.MemoryLimit,
    )
    
    return err
}

func ToggleProblemPublishStatus(problemID int) error {
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback()

    // First check current status
    var isCurrentlyPublished bool
    err = tx.QueryRow("SELECT is_published FROM problems WHERE id = ?", problemID).Scan(&isCurrentlyPublished)
    if err != nil {
        return fmt.Errorf("failed to get current status: %v", err)
    }

    // Prepare update query based on current status
    var query string
    if isCurrentlyPublished {
        // If currently published, unpublish it
        query = `
            UPDATE problems 
            SET is_published = FALSE,
                published_at = NULL
            WHERE id = ?
        `
    } else {
        // If currently unpublished, publish it
        query = `
            UPDATE problems 
            SET is_published = TRUE,
                published_at = CURRENT_TIMESTAMP
            WHERE id = ?
        `
    }

    result, err := tx.Exec(query, problemID)
    if err != nil {
        return fmt.Errorf("failed to update problem status: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error checking affected rows: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("problem not found")
    }

    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }

    return nil
}

func GetSingleProblem(id int) (models.Problem, error) {
    var problem models.Problem
    query := `
        SELECT id, user_id, title, description_path, input_path, output_path, created_at, is_published, time_limit_ms, memory_limit_mb
        FROM problems
        WHERE id = ? LIMIT 1
    `

    err := db.QueryRow(query, id).Scan(
        &problem.ID,
        &problem.UserID,
        &problem.Title,
        &problem.DescriptionPath,
        &problem.InputPath,
        &problem.OutputPath,
        &problem.CreatedTime,
        &problem.IsPublished,
        &problem.TimeLimit,
        &problem.MemoryLimit,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Problem{}, errors.New("problem does not exist")
        }
        return models.Problem{}, err
    }

    return problem, nil
}

func EditProblem(db *sql.DB, user_id, problem_id int, title string, time_limit_ms, memory_limit_mb int) error {
	is_admin, err := GetUserRole(user_id)
	if err != nil {
		return err
	}

	if is_admin == 0 {
		var owner_id int
		err = db.QueryRow("SELECT user_id FROM problems WHERE id = ?", problem_id).Scan(&owner_id)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("problem not found")
			}
			return err
		}
		if owner_id != user_id {
			return errors.New("permission denied: not owner or admin")
		}
	}

	query := `
		UPDATE problems
		SET title = ?, time_limit_ms = ?, memory_limit_mb = ?
		WHERE id = ?
	`
	_, err = db.Exec(query, title, time_limit_ms, memory_limit_mb, problem_id)
	if err != nil {
		return err
	}

	return nil
}

// For users: count only published problems
func GetTotalProblemsCount() (int, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM problems WHERE is_published = true").Scan(&count)
    if err != nil {
        log.Printf("Error counting problems: %v", err)
        return 0, err
    }
    return count, nil
}

// For admins: count all problems
func GetTotalProblemsCountAdmin() (int, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM problems").Scan(&count)
    if err != nil {
        log.Printf("Error counting problems: %v", err)
        return 0, err
    }
    return count, nil
}