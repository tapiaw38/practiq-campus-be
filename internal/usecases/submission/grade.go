package submission

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	GradeUsecase interface {
		Execute(context.Context, string, bool, string, GradeInput) (*GradeOutput, apperrors.ApplicationError)
	}

	gradeUsecase struct {
		contextFactory appcontext.Factory
	}

	GradeInput struct {
		Score        int
		Feedback     string
		RubricScores []RubricScoreInput
	}

	RubricScoreInput struct {
		CriterionID string
		Score       int
		Feedback    string
	}

	GradeOutput struct {
		Data SubmissionData `json:"data"`
	}
)

func NewGradeUsecase(contextFactory appcontext.Factory) GradeUsecase {
	return &gradeUsecase{contextFactory: contextFactory}
}

func (u *gradeUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, submissionID string, input GradeInput) (*GradeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	s, err := app.Repositories.Submission.Get(ctx, submissionID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionGetError, err)
	}
	if s == nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionNotFoundError, nil)
	}

	a, err := app.Repositories.Assignment.Get(ctx, s.AssignmentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentGetError, err)
	}
	if a == nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentNotFoundError, nil)
	}

	if !isSuperAdmin {
		c, err := app.Repositories.Course.Get(ctx, a.CourseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if c == nil || c.OwnerID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

	criteria, err := app.Repositories.Rubric.List(ctx, a.ID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if len(criteria) == 0 {
		if input.Score < 0 || input.Score > a.MaxScore {
			return nil, apperrors.NewBadRequestError("score must be between 0 and the assignment's max score")
		}
		if err := app.Repositories.Submission.Grade(ctx, submissionID, input.Score, input.Feedback); err != nil {
			return nil, apperrors.NewApplicationError(mappings.SubmissionCreateError, err)
		}
	} else {
		if len(input.RubricScores) != len(criteria) {
			return nil, apperrors.NewBadRequestError("a score is required for every rubric criterion")
		}
		byID := make(map[string]domain.RubricCriterion, len(criteria))
		for _, criterion := range criteria {
			byID[criterion.ID] = criterion
		}
		seen := make(map[string]bool, len(criteria))
		scores := make([]domain.RubricScore, 0, len(criteria))
		total := 0
		for _, inputScore := range input.RubricScores {
			criterion, ok := byID[inputScore.CriterionID]
			if !ok || seen[inputScore.CriterionID] || inputScore.Score < 0 || inputScore.Score > criterion.MaxScore {
				return nil, apperrors.NewBadRequestError("invalid rubric score")
			}
			seen[inputScore.CriterionID] = true
			total += inputScore.Score
			scores = append(scores, domain.RubricScore{SubmissionID: submissionID, CriterionID: criterion.ID, Score: inputScore.Score, Feedback: strings.TrimSpace(inputScore.Feedback)})
		}
		if err := app.Repositories.Rubric.Grade(ctx, submissionID, total, input.Feedback, scores); err != nil {
			return nil, apperrors.NewApplicationError(mappings.SubmissionCreateError, err)
		}
	}

	updated, err := app.Repositories.Submission.Get(ctx, submissionID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewInternalError(nil)
	}
	_ = app.Repositories.Notification.Create(ctx, domain.Notification{UserID: updated.UserID, Type: "submission_graded", Title: "Tarea corregida: " + a.Title, Body: "Tu docente publicó una calificación", Data: `{"submission_id":"` + updated.ID + `","assignment_id":"` + a.ID + `"}`})

	data := toSubmissionData(*updated)
	scores, err := app.Repositories.Rubric.Scores(ctx, submissionID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return &GradeOutput{Data: withAttachments(app, withRubricScores(data, scores))}, nil
}
