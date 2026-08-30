package message

import (
	"context"
	"sort"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	SearchRecipientsUsecase interface {
		Execute(context.Context, string, bool, string, string) (*SearchRecipientsOutput, apperrors.ApplicationError)
	}

	searchRecipientsUsecase struct {
		contextFactory appcontext.Factory
	}

	RecipientData struct {
		ID       string `json:"id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}

	SearchRecipientsOutput struct {
		Data []RecipientData `json:"data"`
	}
)

func NewSearchRecipientsUsecase(contextFactory appcontext.Factory) SearchRecipientsUsecase {
	return &searchRecipientsUsecase{contextFactory: contextFactory}
}

// Execute answers "who can I message in this course" — the course's owner
// plus everyone actively enrolled, minus the requester themselves — filtered
// in Go rather than SQL: a course roster is small (a class), so there's no
// real cost, and it avoids a bespoke cross-table search query for what is
// fundamentally list-then-filter.
func (u *searchRecipientsUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID, query string) (*SearchRecipientsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	c, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if c == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}

	if !isSuperAdmin && c.OwnerID != requesterID {
		enrollment, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, courseID, requesterID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
		}
		if enrollment == nil {
			return nil, apperrors.NewForbiddenError()
		}
	}

	candidateIDs := map[string]bool{c.OwnerID: true}
	enrollments, err := app.Repositories.Enrollment.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentListError, err)
	}
	for _, e := range enrollments {
		candidateIDs[e.UserID] = true
	}
	delete(candidateIDs, requesterID)

	q := strings.ToLower(strings.TrimSpace(query))
	results := make([]RecipientData, 0, len(candidateIDs))
	for id := range candidateIDs {
		p, err := app.Repositories.Profile.Get(ctx, id)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
		}
		if p == nil {
			continue
		}
		if q != "" &&
			!strings.Contains(strings.ToLower(p.FullName), q) &&
			!strings.Contains(strings.ToLower(p.Email), q) {
			continue
		}
		results = append(results, RecipientData{ID: p.ID, FullName: p.FullName, Email: p.Email})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].FullName < results[j].FullName })

	return &SearchRecipientsOutput{Data: results}, nil
}
