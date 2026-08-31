package course

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	DuplicateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string) (*DuplicateOutput, apperrors.ApplicationError)
	}

	duplicateUsecase struct {
		contextFactory appcontext.Factory
	}

	DuplicateOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewDuplicateUsecase(contextFactory appcontext.Factory) DuplicateUsecase {
	return &duplicateUsecase{contextFactory: contextFactory}
}

// Execute clones a course's structure — sections, assignments with their
// rubric, quizzes with their questions, and materials — into a new draft
// course the requester owns. It deliberately does not carry over
// enrollments, submissions, quiz attempts, forum threads, calendar events or
// groups: those belong to a specific run of the course, not the course
// itself, so a "new academic year" copy starts empty of them.
//
// ponytail: visible_group_id and unlock_after are cleared rather than
// remapped — a duplicated course has no groups yet, and remapping a
// prerequisite chain across two id spaces is real complexity nobody has
// asked for. Re-wire them by hand after duplicating if you need them; add
// remapping if that becomes a real pain point.
func (u *duplicateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string) (*DuplicateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	original, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if original == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if !isSuperAdmin && original.OwnerID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}

	slug := slugify(original.Title+"-copia") + "-" + randomSuffix()
	newCourseID, err := app.Repositories.Course.Create(ctx, domain.Course{
		OwnerID:     original.OwnerID,
		Title:       original.Title + " (copia)",
		Slug:        slug,
		Description: original.Description,
		Status:      domain.CourseStatusDraft,
		Labels:      original.Labels,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseCreateError, err)
	}

	sections, err := app.Repositories.CourseSection.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	sectionIDMap := make(map[string]string, len(sections))
	for _, s := range sections {
		newID, err := app.Repositories.CourseSection.Create(ctx, domain.CourseSection{
			CourseID:    newCourseID,
			Title:       s.Title,
			Description: s.Description,
			Position:    s.Position,
		})
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		sectionIDMap[s.ID] = newID
	}
	mappedSection := func(oldID *string) *string {
		if oldID == nil {
			return nil
		}
		if newID, ok := sectionIDMap[*oldID]; ok {
			return &newID
		}
		return nil
	}

	assignments, err := app.Repositories.Assignment.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	assignmentIDMap := make(map[string]string, len(assignments))
	for _, a := range assignments {
		newID, err := app.Repositories.Assignment.Create(ctx, domain.Assignment{
			CourseID:    newCourseID,
			SectionID:   mappedSection(a.SectionID),
			Title:       a.Title,
			Description: a.Description,
			MaxScore:    a.MaxScore,
			Weight:      a.Weight,
		})
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		assignmentIDMap[a.ID] = newID

		if criteria, err := app.Repositories.Rubric.List(ctx, a.ID); err == nil && len(criteria) > 0 {
			copied := make([]domain.RubricCriterion, 0, len(criteria))
			for _, c := range criteria {
				copied = append(copied, domain.RubricCriterion{Title: c.Title, Description: c.Description, MaxScore: c.MaxScore})
			}
			_ = app.Repositories.Rubric.Replace(ctx, newID, copied)
		}
	}

	quizzes, err := app.Repositories.Quiz.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	for _, q := range quizzes {
		newID, err := app.Repositories.Quiz.Create(ctx, domain.Quiz{
			CourseID:      newCourseID,
			SectionID:     mappedSection(q.SectionID),
			Title:         q.Title,
			Description:   q.Description,
			TimeLimitSecs: q.TimeLimitSecs,
			MaxAttempts:   q.MaxAttempts,
			Weight:        q.Weight,
		})
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}

		if questions, err := app.Repositories.Quiz.ListQuestions(ctx, q.ID); err == nil && len(questions) > 0 {
			copied := make([]domain.QuizQuestion, 0, len(questions))
			for _, qq := range questions {
				copied = append(copied, domain.QuizQuestion{Type: qq.Type, Statement: qq.Statement, Options: qq.Options, CorrectAnswer: qq.CorrectAnswer, Points: qq.Points})
			}
			_ = app.Repositories.Quiz.ReplaceQuestions(ctx, newID, copied)
		}
	}

	if materials, err := app.Repositories.CourseMaterial.ListByCourse(ctx, courseID); err == nil {
		for _, m := range materials {
			var newAssignmentID *string
			if m.AssignmentID != nil {
				if mapped, ok := assignmentIDMap[*m.AssignmentID]; ok {
					newAssignmentID = &mapped
				}
			}
			_, _ = app.Repositories.CourseMaterial.Create(ctx, domain.CourseMaterial{
				CourseID:     newCourseID,
				AssignmentID: newAssignmentID,
				SectionID:    mappedSection(m.SectionID),
				UploaderID:   requesterID,
				Title:        m.Title,
				Description:  m.Description,
				Kind:         m.Kind,
				URL:          m.URL,
			})
		}
	}

	created, err := app.Repositories.Course.Get(ctx, newCourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &DuplicateOutput{Data: toCourseData(*created)}, nil
}
