package web

import (
	"github.com/gin-gonic/gin"
	profileRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/profile"
	tenantRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/tenant"
	handlerAssignment "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/assignment"
	handlerCalendar "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/calendar_event"
	handlerCourse "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/course"
	handlerGroup "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/course_group"
	handlerMaterial "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/course_material"
	handlerSection "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/course_section"
	handlerDashboard "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/dashboard"
	handlerEnrollment "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/enrollment"
	handlerPost "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/forum_post"
	handlerThread "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/forum_thread"
	handlerMessage "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/message"
	handlerNotification "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/notification"
	handlerPreference "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/preference"
	handlerProfile "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/profile"
	handlerQuiz "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/quiz"
	handlerQuizAttempt "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/quiz_attempt"
	handlerRubric "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/rubric"
	handlerSubmission "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/submission"
	handlerTenant "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/tenant"
	handlerUpload "github.com/tapiaw38/practiq-campus-be/internal/adapters/web/handlers/upload"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/revocation"
	"github.com/tapiaw38/practiq-campus-be/internal/usecases"
)

func RegisterRoutes(app *gin.Engine, uc *usecases.Usecases, profiles profileRepo.Repository, tenants tenantRepo.Repository, practiq practiqapi.Client, revoked *revocation.Checker) {
	api := app.Group("/api")
	api.Use(middlewares.AuthMiddleware(profiles, revoked))

	superAdminOnly := api.Group("/")
	superAdminOnly.Use(middlewares.RequireRoles(middlewares.RoleSuperAdmin))

	api.POST("/profile", handlerProfile.NewSyncHandler(uc.Profile.Sync))
	api.GET("/profile/me", handlerProfile.NewGetMeHandler(uc.Profile.Get))
	api.GET("/tenants/mine", handlerTenant.Mine(uc.Tenant.ListMine))

	superAdminOnly.GET("/tenants", handlerTenant.List(uc.Tenant.List))
	superAdminOnly.GET("/tenants/eligible-schools", handlerTenant.EligibleSchools(uc.Tenant.ListEligibleSchools))
	superAdminOnly.POST("/tenants", handlerTenant.Activate(uc.Tenant.Activate))
	superAdminOnly.PATCH("/tenants/:id/status", handlerTenant.SetStatus(uc.Tenant.SetStatus))

	api.Use(middlewares.RequireTenant(tenants, practiq))

	teacherOnly := api.Group("/")
	teacherOnly.Use(middlewares.RequireTenantStaff())

	api.GET("/school/members", handlerTenant.SchoolMembers(uc.Tenant.ListMembers))
	api.POST("/school/members", handlerTenant.AddSchoolMember(uc.Tenant.AddMember))
	api.DELETE("/school/members/:userID", handlerTenant.RemoveSchoolMember(uc.Tenant.RemoveMember))
	api.GET("/me/preferences/:scope", handlerPreference.NewGetHandler(uc.Preference.Get))
	api.PUT("/me/preferences/:scope", handlerPreference.NewUpdateHandler(uc.Preference.Update))
	api.GET("/notifications", handlerNotification.List(uc.Notification.Manage))
	api.GET("/notifications/summary", handlerNotification.Summary(uc.Notification.Manage))
	api.PUT("/notifications/:id/read", handlerNotification.Read(uc.Notification.Manage))

	superAdminOnly.POST("/users", handlerProfile.NewProvisionStudentHandler(uc.Profile.CreateOrSync))
	superAdminOnly.GET("/users", handlerProfile.NewListStudentsHandler(uc.Profile.ListUsers))
	superAdminOnly.PATCH("/users/:id/block", handlerProfile.NewBlockHandler(uc.Profile.Block))

	teacherOnly.POST("/courses", handlerCourse.NewCreateHandler(uc.Course.Create))
	teacherOnly.POST("/courses/sync-from-practiq", handlerCourse.NewSyncFromPractiqHandler(uc.Course.SyncFromPractiq))
	teacherOnly.POST("/courses/:id/duplicate", handlerCourse.NewDuplicateHandler(uc.Course.Duplicate))
	api.GET("/courses", handlerCourse.NewListHandler(uc.Course.List))
	api.GET("/student/dashboard", handlerDashboard.Student(uc.Dashboard.Student))
	api.GET("/courses/:id", handlerCourse.NewGetHandler(uc.Course.Get))
	teacherOnly.PUT("/courses/:id", handlerCourse.NewUpdateHandler(uc.Course.Update))
	teacherOnly.DELETE("/courses/:id", handlerCourse.NewDeleteHandler(uc.Course.Delete))

	teacherOnly.POST("/courses/:id/enrollments", handlerEnrollment.NewCreateHandler(uc.Enrollment.Create))
	teacherOnly.GET("/courses/:id/enrollments", handlerEnrollment.NewListByCourseHandler(uc.Enrollment.ListByCourse))
	api.GET("/me/enrollments", handlerEnrollment.NewListMineHandler(uc.Enrollment.ListMine))

	teacherOnly.GET("/me/practiq-students", handlerProfile.NewListMyPractiqStudentsHandler(uc.Profile.ListMyPractiqStudents))
	teacherOnly.DELETE("/enrollments/:id", handlerEnrollment.NewDeleteHandler(uc.Enrollment.Delete))

	teacherOnly.POST("/courses/:id/sections", handlerSection.NewCreateHandler(uc.CourseSection.Create))
	api.GET("/courses/:id/sections", handlerSection.NewListHandler(uc.CourseSection.List))
	teacherOnly.PUT("/sections/:id", handlerSection.NewUpdateHandler(uc.CourseSection.Update))
	teacherOnly.DELETE("/sections/:id", handlerSection.NewDeleteHandler(uc.CourseSection.Delete))

	api.GET("/courses/:id/groups", handlerGroup.List(uc.CourseGroup.Manage))
	teacherOnly.POST("/courses/:id/groups", handlerGroup.Create(uc.CourseGroup.Manage))
	teacherOnly.DELETE("/groups/:id", handlerGroup.Delete(uc.CourseGroup.Manage))
	teacherOnly.POST("/groups/:id/members", handlerGroup.AddMember(uc.CourseGroup.Manage))
	teacherOnly.DELETE("/groups/:id/members/:userId", handlerGroup.RemoveMember(uc.CourseGroup.Manage))

	api.POST("/uploads", handlerUpload.NewHandler(uc.Upload.Upload))
	teacherOnly.POST("/courses/:id/materials", handlerMaterial.NewCreateHandler(uc.CourseMaterial.Create))
	api.GET("/courses/:id/materials", handlerMaterial.NewListHandler(uc.CourseMaterial.List))
	teacherOnly.DELETE("/materials/:id", handlerMaterial.NewDeleteHandler(uc.CourseMaterial.Delete))

	teacherOnly.POST("/courses/:id/assignments", handlerAssignment.NewCreateHandler(uc.Assignment.Create))
	api.GET("/courses/:id/assignments", handlerAssignment.NewListHandler(uc.Assignment.List))
	teacherOnly.PUT("/assignments/:id", handlerAssignment.NewUpdateHandler(uc.Assignment.Update))
	teacherOnly.DELETE("/assignments/:id", handlerAssignment.NewDeleteHandler(uc.Assignment.Delete))
	api.GET("/assignments/:id/rubric", handlerRubric.List(uc.Rubric.Manage))
	teacherOnly.PUT("/assignments/:id/rubric", handlerRubric.Replace(uc.Rubric.Manage))

	api.POST("/assignments/:id/submissions", handlerSubmission.NewCreateHandler(uc.Submission.Create))
	api.GET("/assignments/:id/submissions/me", handlerSubmission.NewGetMineHandler(uc.Submission.GetMine))
	teacherOnly.GET("/assignments/:id/submissions", handlerSubmission.NewListByAssignmentHandler(uc.Submission.ListByAssignment))
	teacherOnly.PUT("/submissions/:id/grade", handlerSubmission.NewGradeHandler(uc.Submission.Grade))

	teacherOnly.POST("/courses/:id/quizzes", handlerQuiz.NewCreateHandler(uc.Quiz.Create))
	api.GET("/courses/:id/quizzes", handlerQuiz.NewListHandler(uc.Quiz.List))
	api.GET("/quizzes/:id", handlerQuiz.NewGetHandler(uc.Quiz.Get))
	teacherOnly.PUT("/quizzes/:id", handlerQuiz.NewUpdateHandler(uc.Quiz.Update))
	teacherOnly.DELETE("/quizzes/:id", handlerQuiz.NewDeleteHandler(uc.Quiz.Delete))
	teacherOnly.GET("/quizzes/:id/questions", handlerQuiz.NewListQuestionsHandler(uc.Quiz.Questions))
	teacherOnly.PUT("/quizzes/:id/questions", handlerQuiz.NewReplaceQuestionsHandler(uc.Quiz.Questions))

	api.POST("/quizzes/:id/attempts", handlerQuizAttempt.NewStartHandler(uc.QuizAttempt.Start))
	api.GET("/quizzes/:id/attempts/mine", handlerQuizAttempt.NewListMineHandler(uc.QuizAttempt.ListMine))
	teacherOnly.GET("/quizzes/:id/attempts", handlerQuizAttempt.NewListByQuizHandler(uc.QuizAttempt.ListByQuiz))
	api.PUT("/attempts/:id/draft", handlerQuizAttempt.NewSaveDraftHandler(uc.QuizAttempt.SaveDraft))
	api.POST("/attempts/:id/submit", handlerQuizAttempt.NewSubmitHandler(uc.QuizAttempt.Submit))
	api.GET("/attempts/:id", handlerQuizAttempt.NewGetHandler(uc.QuizAttempt.Get))

	teacherOnly.POST("/courses/:id/forum-threads", handlerThread.NewCreateHandler(uc.ForumThread.Create))
	teacherOnly.PUT("/forum-threads/:id", handlerThread.NewUpdateHandler(uc.ForumThread.Update))
	api.GET("/courses/:id/forum-threads", handlerThread.NewListHandler(uc.ForumThread.List))
	api.POST("/forum-threads/:id/posts", handlerPost.NewCreateHandler(uc.ForumPost.Create))
	api.GET("/forum-threads/:id/posts", handlerPost.NewListHandler(uc.ForumPost.List))
	api.PUT("/forum-posts/:id", handlerPost.NewUpdateHandler(uc.ForumPost.Update))
	api.DELETE("/forum-posts/:id", handlerPost.NewDeleteHandler(uc.ForumPost.Delete))

	teacherOnly.POST("/calendar/events", handlerCalendar.NewCreateHandler(uc.CalendarEvent.Create))
	teacherOnly.PUT("/calendar/events/:id", handlerCalendar.NewUpdateHandler(uc.CalendarEvent.Update))
	teacherOnly.DELETE("/calendar/events/:id", handlerCalendar.NewDeleteHandler(uc.CalendarEvent.Delete))
	api.GET("/me/calendar", handlerCalendar.NewListMineHandler(uc.CalendarEvent.ListMine))

	api.POST("/messages", handlerMessage.NewSendHandler(uc.Message.Send))
	api.GET("/courses/:id/messageable-users", handlerMessage.NewSearchRecipientsHandler(uc.Message.SearchRecipients))
	teacherOnly.POST("/courses/:id/messages/broadcast", handlerMessage.NewBroadcastHandler(uc.Message.Broadcast))
	api.GET("/conversations", handlerMessage.NewListConversationsHandler(uc.Message.ListConversations))
	api.GET("/conversations/:id/messages", handlerMessage.NewListMessagesHandler(uc.Message.ListMessages))
	api.POST("/conversations/:id/messages", handlerMessage.NewReplyHandler(uc.Message.Reply))
	api.POST("/conversations/:id/read", handlerMessage.NewMarkReadHandler(uc.Message.MarkRead))
}
