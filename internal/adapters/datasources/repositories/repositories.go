package repositories

import (
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course_section"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/enrollment"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/profile"
)

type Repositories struct {
	Profile       profile.Repository
	Course        course.Repository
	CourseSection course_section.Repository
	Enrollment    enrollment.Repository
}

type Factory func() *Repositories

func NewFactory(ds *datasources.Datasources) func() *Repositories {
	return func() *Repositories {
		return &Repositories{
			Profile:       profile.NewRepository(ds.DB),
			Course:        course.NewRepository(ds.DB),
			CourseSection: course_section.NewRepository(ds.DB),
			Enrollment:    enrollment.NewRepository(ds.DB),
		}
	}
}
