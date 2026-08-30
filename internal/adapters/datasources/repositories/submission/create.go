package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, s domain.Submission) (string, error) {
	query := `
		INSERT INTO submissions (assignment_id, user_id, content, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, s.AssignmentID, s.UserID, s.Content, domain.SubmissionStatusSubmitted).Scan(&id)
	return id, err
}
