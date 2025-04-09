package database

import (
	// "database/sql"
	// "errors"

	"github.com/placeHolder143032/CodeChallengeHub/models"

	_ "github.com/go-sql-driver/mysql"
)

func SubmitCode(submission models.Submission) error {
	insertQuery := `
		INSERT INTO submissions (user_id, problem_id, code_path, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(insertQuery, submission.User_id, submission.Problem_id, submission.Code_path, submission.Created_at)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSubmission(submission models.Submission) error {
	query := `
		UPDATE submissions
		SET state = ?, runtime_ms = ?, memory_used = ?, error_message = ?
		WHERE id = ?
	`

	_, err := db.Exec(query, submission.State, submission.Runtime_ms, submission.Memory_used, submission.Error_message, submission.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllSubmissionsByUser(user_id int) ([]models.Submission, error) {
	query := `
		SELECT problem_id, code_path, state, created_at, runtime_ms, memory_used, error_message
		FROM submissions
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []models.Submission
	for rows.Next() {
		var sub models.Submission
		err := rows.Scan(&sub.Problem_id, &sub.Code_path, &sub.State, &sub.Created_at, &sub.Runtime_ms, &sub.Memory_used, &sub.Error_message)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, sub)
	}

	return submissions, nil
}
