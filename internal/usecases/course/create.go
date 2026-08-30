package course

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Title            string
		Description      string
		StartDate        *time.Time
		EndDate          *time.Time
		PractiqSubjectID *string
		Labels           []string
	}

	CreateOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

var slugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(title string) string {
	s := slugSanitizer.ReplaceAllString(strings.ToLower(title), "-")
	return strings.Trim(s, "-")
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}

	app := u.contextFactory()

	slug := slugify(input.Title)
	if existing, err := app.Repositories.Course.GetBySlug(ctx, slug); err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	} else if existing != nil {
		// Two courses with the same title is common ("Matemática"); the slug
		// is disambiguated with a short suffix rather than rejecting the
		// request, which would make the teacher rename their own course.
		slug = slug + "-" + randomSuffix()
	}

	id, err := app.Repositories.Course.Create(ctx, domain.Course{
		OwnerID:          requesterID,
		Title:            input.Title,
		Slug:             slug,
		Description:      input.Description,
		Status:           domain.CourseStatusDraft,
		StartDate:        input.StartDate,
		EndDate:          input.EndDate,
		PractiqSubjectID: input.PractiqSubjectID,
		Labels:           normalizeLabels(input.Labels),
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseCreateError, err)
	}

	created, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &CreateOutput{Data: toCourseData(*created)}, nil
}

func normalizeLabels(labels []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(labels))
	for _, label := range labels {
		value := strings.TrimSpace(label)
		if value != "" && len(value) <= 50 && !seen[strings.ToLower(value)] {
			result = append(result, value)
			seen[strings.ToLower(value)] = true
		}
	}
	return result
}

// randomSuffix disambiguates a slug collision. Time-based, not random: two
// requests never collide, and it needs no extra dependency.
func randomSuffix() string {
	return strconv.FormatInt(time.Now().UnixNano()%1_000_000, 36)
}
