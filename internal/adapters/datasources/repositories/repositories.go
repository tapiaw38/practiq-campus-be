package repositories

import (
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/assignment"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/calendar_event"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/conversation"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course_group"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course_material"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course_section"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/enrollment"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/forum_post"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/forum_thread"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/notification"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/preference"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/profile"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/quiz"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/quiz_attempt"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/rubric"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/submission"
)

type Repositories struct {
	Profile        profile.Repository
	Preference     preference.Repository
	Course         course.Repository
	CourseSection  course_section.Repository
	CourseGroup    course_group.Repository
	CourseMaterial course_material.Repository
	Enrollment     enrollment.Repository
	Assignment     assignment.Repository
	Submission     submission.Repository
	Rubric         rubric.Repository
	Quiz           quiz.Repository
	QuizAttempt    quiz_attempt.Repository
	ForumThread    forum_thread.Repository
	ForumPost      forum_post.Repository
	CalendarEvent  calendar_event.Repository
	Conversation   conversation.Repository
	Notification   notification.Repository
}

type Factory func() *Repositories

func NewFactory(ds *datasources.Datasources) func() *Repositories {
	return func() *Repositories {
		return &Repositories{
			Profile:        profile.NewRepository(ds.DB),
			Preference:     preference.NewRepository(ds.DB),
			Course:         course.NewRepository(ds.DB),
			CourseSection:  course_section.NewRepository(ds.DB),
			CourseGroup:    course_group.NewRepository(ds.DB),
			CourseMaterial: course_material.NewRepository(ds.DB),
			Enrollment:     enrollment.NewRepository(ds.DB),
			Assignment:     assignment.NewRepository(ds.DB),
			Submission:     submission.NewRepository(ds.DB),
			Rubric:         rubric.NewRepository(ds.DB),
			Quiz:           quiz.NewRepository(ds.DB),
			QuizAttempt:    quiz_attempt.NewRepository(ds.DB),
			ForumThread:    forum_thread.NewRepository(ds.DB),
			ForumPost:      forum_post.NewRepository(ds.DB),
			CalendarEvent:  calendar_event.NewRepository(ds.DB),
			Conversation:   conversation.NewRepository(ds.DB),
			Notification:   notification.NewRepository(ds.DB),
		}
	}
}
