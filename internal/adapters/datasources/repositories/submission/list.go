package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) ListByAssignment(ctx context.Context, assignmentID string) ([]domain.Submission, error) {
	query := `SELECT ` + selectSubmissionColumns + ` FROM submissions WHERE assignment_id = $1 ORDER BY submitted_at ASC`
	rows, err := r.db.QueryContext(ctx, query, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []domain.Submission
	for rows.Next() {
		var s domain.Submission
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.UserID, &s.Content, &s.Status, &s.Score, &s.Feedback, &s.SubmittedAt, &s.GradedAt); err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	return submissions, rows.Err()
}
