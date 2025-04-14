package database

import (
	"database/sql"
	"errors"

	"github.com/placeHolder143032/CodeChallengeHub/models"

	_ "github.com/go-sql-driver/mysql"
)

// this one is for admins since it doesnt check if its published or not
func GetProblemsPageAdmin(m, n int) ([]models.Problem, error) {
	var problems []models.Problem
	query := `
		SELECT id, title
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
		if err := rows.Scan(&p.ID, &p.Title); err != nil {
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
        SELECT id, title, time_limit_ms, memory_limit_mb, created_at, is_published
        FROM problems
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?;
    `
    offset := (m - 1) * n

    rows, err := db.Query(query, n, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var p models.Problem
        if err := rows.Scan(&p.ID, &p.Title, &p.TimeLimit, &p.MemoryLimit, &p.CreatedTime, &p.IsPublished); err != nil {
            return nil, err
        }
        problems = append(problems, p)
    }

    if err := rows.Err(); err != nil {
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

func PublishProblem(id int) error {
	query := `
	UPDATE problems SET is_published = NOT is_published
	WHERE id = ? LIMIT 1
	`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetSingleProblem(id int) (models.Problem, error) {
	var problem models.Problem
	query := `
	SELECT title, description_path, input_path, output_path, time_limit_ms, memory_limit_mb FROM users 
	WHERE id = ? LIMIT 1
	`

	err := db.QueryRow(query, id).Scan(&problem.Title, problem.DescriptionPath, problem.InputPath, problem.OutputPath, problem.CreatedTime, problem.TimeLimit, problem.MemoryLimit)
	if err != nil && err != sql.ErrNoRows {
		return models.Problem{}, err
	}

	if err == sql.ErrNoRows {
		return models.Problem{}, errors.New("user does not exist")
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

func GetTotalProblemsCount() (int, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM problems").Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}