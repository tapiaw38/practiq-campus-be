package dashboard

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

type StudentUsecase interface {
	Execute(context.Context, string) (*StudentOutput, apperrors.ApplicationError)
}
type studentUsecase struct{ contextFactory appcontext.Factory }
type CourseData struct {
	CourseID        string  `json:"course_id"`
	ProgressPercent int     `json:"progress_percent"`
	PendingItems    int     `json:"pending_items"`
	NextDueAt       *string `json:"next_due_at"`
	NextDueTitle    string  `json:"next_due_title"`
}
type StudentOutput struct {
	Data []CourseData `json:"data"`
}

func NewStudentUsecase(f appcontext.Factory) StudentUsecase {
	return &studentUsecase{contextFactory: f}
}
func (u *studentUsecase) Execute(ctx context.Context, userID string) (*StudentOutput, apperrors.ApplicationError) {
	items, err := u.contextFactory().Repositories.Dashboard.StudentCourses(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	out := make([]CourseData, 0, len(items))
	for _, item := range items {
		progress := 0
		if item.TotalItems > 0 {
			progress = item.CompletedItems * 100 / item.TotalItems
		}
		var due *string
		if item.NextDueAt != nil {
			value := item.NextDueAt.Format("2006-01-02T15:04:05Z")
			due = &value
		}
		out = append(out, CourseData{CourseID: item.CourseID, ProgressPercent: progress, PendingItems: item.PendingItems, NextDueAt: due, NextDueTitle: item.NextDueTitle})
	}
	return &StudentOutput{Data: out}, nil
}
