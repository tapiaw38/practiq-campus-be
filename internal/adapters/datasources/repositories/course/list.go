package course

import (
	"context"
	"strconv"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) List(ctx context.Context, filter ListFilter) ([]domain.Course, error) {
	query := "SELECT " + selectCourseColumns + " FROM courses c"
	args := []any{}
	where := []string{}

	if filter.EnrolledUserID != "" {
		query += " JOIN enrollments e ON e.course_id = c.id"
		args = append(args, filter.EnrolledUserID)
		where = append(where, "e.user_id = $"+strconv.Itoa(len(args))+" AND e.status = 'active'")
	}
	if filter.OwnerID != "" {
		args = append(args, filter.OwnerID)
		where = append(where, "c.owner_id = $"+strconv.Itoa(len(args)))
	}
	if filter.PublishedOnly {
		where = append(where, "c.status = 'published'")
	}

	for i, cond := range where {
		if i == 0 {
			query += " WHERE "
		} else {
			query += " AND "
		}
		query += cond
	}
	query += " ORDER BY c.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(
			&c.ID, &c.OwnerID, &c.Title, &c.Slug, &c.Description, &c.Status,
			&c.StartDate, &c.EndDate, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}
