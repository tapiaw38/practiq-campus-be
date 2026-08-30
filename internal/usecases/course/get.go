package course

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(context.Context, string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

// Execute is intentionally not ownership/enrollment-gated in Phase 1: there
// is no course catalog or "preview" concept yet, and every caller reaching
// this either owns the course or was enrolled by its owner. Revisit once a
// public catalog exists.
func (u *getUsecase) Execute(ctx context.Context, courseID string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	c, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if c == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}

	return &GetOutput{Data: toCourseData(*c)}, nil
}
