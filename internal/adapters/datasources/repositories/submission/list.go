package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) ListByAssignment(ctx context.Context, assignmentID string) ([]domain.Submission, error) {
	query, args, err := tenantChain(ctx).Where("s.assignment_id = ?", assignmentID).
		SQL(selectQualifiedSubmissionColumns, "submissions s")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY s.submitted_at ASC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []domain.Submission
	for rows.Next() {
		s, err := scanSubmission(rows)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, *s)
	}
	return submissions, rows.Err()
}
