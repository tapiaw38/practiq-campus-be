package submission

import (
	"context"
	"errors"
	"strings"

	rubricrepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/rubric"
	submissionrepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/submission"
	"github.com/tapiaw38/practiq-campus-be/internal/usecases/campusaccess"

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
		Score    int
		Feedback string
		// Version is the submission the teacher actually reviewed. Zero means
		// the client did not send one, which is accepted as "whatever is
		// current" so older clients keep working — they simply do not get the
		// protection.
		Version      int
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
		if c == nil || !campusaccess.CanManageCourse(ctx, c.OwnerID, requesterID, false) {
			return nil, apperrors.NewForbiddenError()
		}
	}

	// A client that sends no version is grading whatever is current; one that
	// does is grading the version it read, and loses if that moved.
	expectedVersion := input.Version
	if expectedVersion == 0 {
		expectedVersion = s.Version
	}

	criteria, err := app.Repositories.Rubric.List(ctx, a.ID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if len(criteria) == 0 {
		if input.Score < 0 || input.Score > a.MaxScore {
			return nil, apperrors.NewBadRequestError("score must be between 0 and the assignment's max score")
		}
		if err := app.Repositories.Submission.Grade(ctx, submissionID, input.Score, input.Feedback, expectedVersion); err != nil {
			if errors.Is(err, submissionrepo.ErrSubmissionChanged) {
				return nil, apperrors.NewBadRequestError("el alumno reenvió la entrega mientras la corregías: volvé a abrirla para revisar la versión nueva")
			}
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
		// Each criterion is checked against its own maximum, but the sum was
		// never checked against the assignment's. The two are kept equal when
		// either is edited — now in both directions — yet an assignment whose
		// maximum was lowered before that rule existed can still hold a
		// rubric totalling more, and recording 100 out of 10 would make the
		// grade meaningless.
		if total > a.MaxScore {
			return nil, apperrors.NewBadRequestError("la rúbrica suma más que el puntaje máximo de la tarea: revisala antes de corregir")
		}
		if err := app.Repositories.Rubric.Grade(ctx, submissionID, total, input.Feedback, scores, expectedVersion); err != nil {
			if errors.Is(err, rubricrepo.ErrSubmissionChanged) {
				return nil, apperrors.NewBadRequestError("el alumno reenvió la entrega mientras la corregías: volvé a abrirla para revisar la versión nueva")
			}
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
