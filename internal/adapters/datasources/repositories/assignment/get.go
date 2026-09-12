package assignment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Qualified for the tenant join, which brings courses into scope.
const selectQualifiedAssignmentColumns = `
	a.id, a.course_id, a.section_id, a.title, a.description, a.due_at, a.max_score, a.created_at, a.updated_at,
	a.weight, a.visible_group_id, a.unlock_after_type, a.unlock_after_id
`

func scanAssignment(row interface{ Scan(...any) error }) (*domain.Assignment, error) {
	var a domain.Assignment
	err := row.Scan(&a.ID, &a.CourseID, &a.SectionID, &a.Title, &a.Description, &a.DueAt, &a.MaxScore,
		&a.CreatedAt, &a.UpdatedAt, &a.Weight, &a.VisibleGroupID, &a.UnlockAfterType, &a.UnlockAfterID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Get reads one assignment of the tenant in context. Assignments carry no
// tenant of their own: they belong to one through their course, and the join
// is what keeps an id from another institution from resolving at all.
func (r *repository) Get(ctx context.Context, id string) (*domain.Assignment, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "a.course_id").
		Where("a.id = ?", id).
		SQL(selectQualifiedAssignmentColumns, "assignments a")
	if err != nil {
		return nil, err
	}
	return scanAssignment(r.db.QueryRowContext(ctx, query, args...))
}
