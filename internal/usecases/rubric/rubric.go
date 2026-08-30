package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

type Criterion struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score"`
}

type Usecase interface {
	List(context.Context, string, string, bool) ([]domain.RubricCriterion, apperrors.ApplicationError)
	Replace(context.Context, string, string, bool, []Criterion) apperrors.ApplicationError
}

type usecase struct{ f appcontext.Factory }

func New(f appcontext.Factory) Usecase { return &usecase{f} }
