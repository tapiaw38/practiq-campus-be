package usecases

import (
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	ucCourse "github.com/tapiaw38/practiq-campus-be/internal/usecases/course"
	ucEnrollment "github.com/tapiaw38/practiq-campus-be/internal/usecases/enrollment"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
)

type ProfileUsecases struct {
	Sync ucProfile.SyncUsecase
	Get  ucProfile.GetUsecase
}

type CourseUsecases struct {
	Create ucCourse.CreateUsecase
	List   ucCourse.ListUsecase
	Get    ucCourse.GetUsecase
	Update ucCourse.UpdateUsecase
}

type EnrollmentUsecases struct {
	Create       ucEnrollment.CreateUsecase
	ListByCourse ucEnrollment.ListByCourseUsecase
	ListMine     ucEnrollment.ListMineUsecase
	Delete       ucEnrollment.DeleteUsecase
}

type Usecases struct {
	Profile    ProfileUsecases
	Course     CourseUsecases
	Enrollment EnrollmentUsecases
}

func NewUsecases(contextFactory appcontext.Factory) *Usecases {
	return &Usecases{
		Profile: ProfileUsecases{
			Sync: ucProfile.NewSyncUsecase(contextFactory),
			Get:  ucProfile.NewGetUsecase(contextFactory),
		},
		Course: CourseUsecases{
			Create: ucCourse.NewCreateUsecase(contextFactory),
			List:   ucCourse.NewListUsecase(contextFactory),
			Get:    ucCourse.NewGetUsecase(contextFactory),
			Update: ucCourse.NewUpdateUsecase(contextFactory),
		},
		Enrollment: EnrollmentUsecases{
			Create:       ucEnrollment.NewCreateUsecase(contextFactory),
			ListByCourse: ucEnrollment.NewListByCourseUsecase(contextFactory),
			ListMine:     ucEnrollment.NewListMineUsecase(contextFactory),
			Delete:       ucEnrollment.NewDeleteUsecase(contextFactory),
		},
	}
}
