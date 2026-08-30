package submission

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectSubmissionColumns = `
	id, assignment_id, user_id, content, status, score, feedback, submitted_at, graded_at
`

func scanSubmission(row *sql.Row) (*domain.Submission, error) {
	var s domain.Submission
	err := row.Scan(&s.ID, &s.AssignmentID, &s.UserID, &s.Content, &s.Status, &s.Score, &s.Feedback, &s.SubmittedAt, &s.GradedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Submission, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectSubmissionColumns+" FROM submissions WHERE id = $1", id)
	return scanSubmission(row)
}

func (r *repository) GetByAssignmentAndUser(ctx context.Context, assignmentID, userID string) (*domain.Submission, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectSubmissionColumns+" FROM submissions WHERE assignment_id = $1 AND user_id = $2", assignmentID, userID)
	return scanSubmission(row)
}
