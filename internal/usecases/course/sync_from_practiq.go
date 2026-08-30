package course

import (
	"context"

	courseRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	SyncFromPractiqUsecase interface {
		Execute(context.Context, string, string) (*SyncFromPractiqOutput, apperrors.ApplicationError)
	}

	syncFromPractiqUsecase struct {
		contextFactory appcontext.Factory
	}

	SyncFromPractiqOutput struct {
		Data []CourseData `json:"data"`
	}
)

func NewSyncFromPractiqUsecase(contextFactory appcontext.Factory) SyncFromPractiqUsecase {
	return &syncFromPractiqUsecase{contextFactory: contextFactory}
}

// Execute pulls the requesting teacher's own practiq-be subjects and creates
// a matching Campus course for every one that doesn't have a course linked
// to it yet — the whole point being the teacher never has to know or type a
// practiq subject id themselves. Already-linked subjects are skipped, so
// calling this repeatedly is safe.
func (u *syncFromPractiqUsecase) Execute(ctx context.Context, requesterID, bearerToken string) (*SyncFromPractiqOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	subjects, err := app.Integrations.PractiqAPI.ListSubjects(ctx, bearerToken)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PractiqLookupError, err)
	}

	existing, err := app.Repositories.Course.List(ctx, courseRepo.ListFilter{OwnerID: requesterID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
	}

	linked := make(map[string]bool, len(existing))
	for _, c := range existing {
		if c.PractiqSubjectID != nil {
			linked[*c.PractiqSubjectID] = true
		}
	}

	for _, subject := range subjects {
		if subject.CreatedBy != requesterID || linked[subject.ID] {
			continue
		}

		slug := slugify(subject.Name)
		if bySlug, err := app.Repositories.Course.GetBySlug(ctx, slug); err == nil && bySlug != nil {
			slug = slug + "-" + randomSuffix()
		}

		subjectID := subject.ID
		if _, err := app.Repositories.Course.Create(ctx, domain.Course{
			OwnerID:          requesterID,
			Title:            subject.Name,
			Slug:             slug,
			Description:      subject.Description,
			Status:           domain.CourseStatusDraft,
			PractiqSubjectID: &subjectID,
		}); err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseCreateError, err)
		}
	}

	updated, err := app.Repositories.Course.List(ctx, courseRepo.ListFilter{OwnerID: requesterID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
	}

	return &SyncFromPractiqOutput{Data: toCourseDataList(updated)}, nil
}
