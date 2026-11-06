package submission

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

const selectQualifiedSubmissionColumns = `
	s.id, s.assignment_id, s.user_id, s.content, s.status, s.version, s.score, s.feedback, s.submitted_at, s.graded_at
`

func tenantChain(ctx context.Context) *tenantcontext.Query {
	return tenantcontext.NewQuery(ctx).
		Join("assignments", "a", "s.assignment_id").
		Through("courses", "c", "a.course_id")
}

func scanSubmission(row interface{ Scan(...any) error }) (*domain.Submission, error) {
	var s domain.Submission
	err := row.Scan(&s.ID, &s.AssignmentID, &s.UserID, &s.Content, &s.Status, &s.Version, &s.Score,
		&s.Feedback, &s.SubmittedAt, &s.GradedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Submission, error) {
	query, args, err := tenantChain(ctx).Where("s.id = ?", id).
		SQL(selectQualifiedSubmissionColumns, "submissions s")
	if err != nil {
		return nil, err
	}
	return scanSubmission(r.db.QueryRowContext(ctx, query, args...))
}

func (r *repository) GetByAssignmentAndUser(ctx context.Context, assignmentID, userID string) (*domain.Submission, error) {
	query, args, err := tenantChain(ctx).
		Where("s.assignment_id = ?", assignmentID).
		Where("s.user_id = ?", userID).
		SQL(selectQualifiedSubmissionColumns, "submissions s")
	if err != nil {
		return nil, err
	}
	return scanSubmission(r.db.QueryRowContext(ctx, query, args...))
}
