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
		FROM questions
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
		SELECT id, title
		FROM questions
		WHERE is_published = true
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

// created at should be from the other side for less inconsistency
func AddProblem(user_id int, problem models.Problem) error {
	insertQuery := `
	INSERT INTO problems (userID, title, descriptionPath, inputPath, outputPath, createdTime, timeLimit, memoryLimit)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(insertQuery, problem.UserID, problem.Title, problem.DescriptionPath, problem.InputPath, problem.OutputPath, problem.CreatedTime, problem.TimeLimit, problem.MemoryLimit)

	if err != nil {
		return err
	}

	return nil
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
