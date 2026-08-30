package course_material

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, bool, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		SectionID   *string
		Title       string
		Description string
		Kind        string
		URL         string
	}

	CreateOutput struct {
		Data MaterialData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}
	if strings.TrimSpace(input.URL) == "" {
		return nil, apperrors.NewBadRequestError("url is required")
	}
	if input.Kind != domain.MaterialKindFile && input.Kind != domain.MaterialKindLink {
		return nil, apperrors.NewBadRequestError("kind must be file or link")
	}
	if input.Kind == domain.MaterialKindLink {
		parsed, err := url.ParseRequestURI(input.URL)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
			return nil, apperrors.NewBadRequestError("url must be an http or https address")
		}
	}

	app := u.contextFactory()

	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}
	if input.SectionID != nil {
		section, err := app.Repositories.CourseSection.Get(ctx, *input.SectionID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.SectionGetError, err)
		}
		if section == nil || section.CourseID != courseID {
			return nil, apperrors.NewBadRequestError("section does not belong to this course")
		}
	}

	// A canonical bucket URL is guessable and readable from other courses'
	// material responses. Without this check a teacher could paste someone
	// else's object into their own course, and withViewURL would then presign
	// it for every reader of that course.
	if input.Kind == domain.MaterialKindFile {
		if app.Storage == nil || !app.Storage.OwnsFileURL(input.URL, MaterialsFolder, requesterID) {
			return nil, apperrors.NewApplicationError(mappings.MaterialFileNotOwnedError,
				errors.New("the file does not belong to this teacher"))
		}
	}

	id, err := app.Repositories.CourseMaterial.Create(ctx, domain.CourseMaterial{
		CourseID:    courseID,
		SectionID:   input.SectionID,
		UploaderID:  requesterID,
		Title:       input.Title,
		Description: input.Description,
		Kind:        input.Kind,
		URL:         input.URL,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialCreateError, err)
	}

	created, err := app.Repositories.CourseMaterial.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &CreateOutput{Data: withViewURL(app, toMaterialData(*created))}, nil
}
