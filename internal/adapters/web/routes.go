package web

import (
	"github.com/gin-gonic/gin"
	handlerCourse "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/course"
	handlerEnrollment "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/enrollment"
	handlerProfile "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/profile"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-campus-be/internal/usecases"
)

func RegisterRoutes(app *gin.Engine, uc *usecases.Usecases) {
	api := app.Group("/api")
	api.Use(middlewares.AuthMiddleware())

	teacherOnly := api.Group("/")
	teacherOnly.Use(middlewares.RequireRoles(middlewares.RoleTeacher, middlewares.RoleSuperAdmin))

	// Profile — any authenticated user.
	api.POST("/profile", handlerProfile.NewSyncHandler(uc.Profile.Sync))
	api.GET("/profile/me", handlerProfile.NewGetMeHandler(uc.Profile.Get))

	// Courses.
	teacherOnly.POST("/courses", handlerCourse.NewCreateHandler(uc.Course.Create))
	api.GET("/courses", handlerCourse.NewListHandler(uc.Course.List))
	api.GET("/courses/:id", handlerCourse.NewGetHandler(uc.Course.Get))
	teacherOnly.PUT("/courses/:id", handlerCourse.NewUpdateHandler(uc.Course.Update))

	// Enrollments.
	teacherOnly.POST("/courses/:id/enrollments", handlerEnrollment.NewCreateHandler(uc.Enrollment.Create))
	teacherOnly.GET("/courses/:id/enrollments", handlerEnrollment.NewListByCourseHandler(uc.Enrollment.ListByCourse))
	api.GET("/me/enrollments", handlerEnrollment.NewListMineHandler(uc.Enrollment.ListMine))
	teacherOnly.DELETE("/enrollments/:id", handlerEnrollment.NewDeleteHandler(uc.Enrollment.Delete))
}
