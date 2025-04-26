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

	_, err := db.Exec(insertQuery, submission.UserId, submission.ProblemId, submission.CodePath, submission.CreatedAt)
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
		SELECT s.id, s.problem_id, p.title, s.code_path, s.state, s.created_at, s.runtime_ms, s.memory_used, s.error_message
		FROM submissions s
		JOIN problems p ON s.problem_id = p.id
		WHERE s.user_id = ?
		ORDER BY s.created_at DESC
	`

	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []models.Submission
	for rows.Next() {
		var sub models.Submission
		err := rows.Scan(&sub.ID, &sub.ProblemId, &sub.ProblemTitle, &sub.CodePath, &sub.State, &sub.CreatedAt, &sub.Runtime_ms, &sub.Memory_used, &sub.Error_message)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, sub)
	}

	return submissions, nil
}


func GetTotalSubmissionsCount(userID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM submissions WHERE user_id = ?"
	
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

func GetStatusFromState(state int) string {
	switch state {
	case 0:
		return "Pending"
	case 1:
		return "OK"
	case 2:
		return "Compile Error"
	case 3:
		return "Wrong Answer"
	case 4:
		return "Memory Limit"
	case 5:
		return "Time Limit"
	case 6:
		return "Runtime Error"
	default:
		return "Unknown"
	}
}

func UpdateSubmissionState(submissionID int64, state int) error {
    query := "UPDATE submissions SET state = ? WHERE id = ?"
    _, err := db.Exec(query, state, submissionID)
    return err
}