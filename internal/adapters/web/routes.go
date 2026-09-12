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
	"github.com/tapiaw38/practiq-campus-be/internal/usecases"
)

func RegisterRoutes(app *gin.Engine, uc *usecases.Usecases, profiles profileRepo.Repository, tenants tenantRepo.Repository, practiq practiqapi.Client) {
	api := app.Group("/api")
	api.Use(middlewares.AuthMiddleware(profiles))

	// These platform administration routes intentionally exist outside a
	// tenant: they are how a superadmin creates and selects tenants.
	superAdminOnly := api.Group("/")
	superAdminOnly.Use(middlewares.RequireRoles(middlewares.RoleSuperAdmin))

	// Profile — any authenticated user.
	api.POST("/profile", handlerProfile.NewSyncHandler(uc.Profile.Sync))
	api.GET("/profile/me", handlerProfile.NewGetMeHandler(uc.Profile.Get))
	api.GET("/tenants/mine", handlerTenant.Mine(tenants, practiq))

	// Campus institutions are enabled by the platform superadmin, never by an
	// institution itself. These are registered on a group created before
	// RequireTenant, so they stay reachable without one: they are how a tenant
	// comes to exist, and an operator has none selected while doing it.
	superAdminOnly.GET("/tenants", handlerTenant.List(tenants, practiq))
	superAdminOnly.GET("/tenants/eligible-schools", handlerTenant.EligibleSchools(tenants, practiq))
	superAdminOnly.POST("/tenants", handlerTenant.Activate(tenants, practiq))
	superAdminOnly.PATCH("/tenants/:id/status", handlerTenant.SetStatus(tenants, practiq))

	// Every product route below runs inside a validated Campus institution.
	// The three bootstrap routes above deliberately do not: they establish the
	// local profile and let a user discover which tenant to select.
	api.Use(middlewares.RequireTenant(tenants, practiq))
	// This group must be created after RequireTenant. Gin copies parent
	// handlers at Group creation; creating it above would let a global auth
	// role decide access before Campus knew which institution was selected.
	teacherOnly := api.Group("/")
	teacherOnly.Use(middlewares.RequireTenantStaff())
	// Institution administration is tenant-scoped. A school admin can only see
	// their selected school's members; platform superadmin can operate any
	// selected Campus institution.
	api.GET("/school/members", handlerTenant.SchoolMembers(tenants, practiq))
	api.POST("/school/members", handlerTenant.AddSchoolMember(tenants, practiq))
	api.DELETE("/school/members/:userID", handlerTenant.RemoveSchoolMember(tenants, practiq))
	api.GET("/me/preferences/:scope", handlerPreference.NewGetHandler(uc.Preference.Get))
	api.PUT("/me/preferences/:scope", handlerPreference.NewUpdateHandler(uc.Preference.Update))
	api.GET("/notifications", handlerNotification.List(uc.Notification.Manage))
	api.PUT("/notifications/:id/read", handlerNotification.Read(uc.Notification.Manage))

	// Users — superadmin only: creating (and listing) accounts on behalf of
	// someone else is an administrative action, not something any teacher
	// can do.
	superAdminOnly.POST("/users", handlerProfile.NewProvisionStudentHandler(uc.Profile.CreateOrSync))
	superAdminOnly.GET("/users", handlerProfile.NewListStudentsHandler(uc.Profile.ListUsers))
	superAdminOnly.PATCH("/users/:id/block", handlerProfile.NewBlockHandler(uc.Profile.Block))

	// Courses.
	teacherOnly.POST("/courses", handlerCourse.NewCreateHandler(uc.Course.Create))
	teacherOnly.POST("/courses/sync-from-practiq", handlerCourse.NewSyncFromPractiqHandler(uc.Course.SyncFromPractiq))
	teacherOnly.POST("/courses/:id/duplicate", handlerCourse.NewDuplicateHandler(uc.Course.Duplicate))
	api.GET("/courses", handlerCourse.NewListHandler(uc.Course.List))
	api.GET("/courses/:id", handlerCourse.NewGetHandler(uc.Course.Get))
	teacherOnly.PUT("/courses/:id", handlerCourse.NewUpdateHandler(uc.Course.Update))
	teacherOnly.DELETE("/courses/:id", handlerCourse.NewDeleteHandler(uc.Course.Delete))

	// Enrollments.
	teacherOnly.POST("/courses/:id/enrollments", handlerEnrollment.NewCreateHandler(uc.Enrollment.Create))
	teacherOnly.GET("/courses/:id/enrollments", handlerEnrollment.NewListByCourseHandler(uc.Enrollment.ListByCourse))
	api.GET("/me/enrollments", handlerEnrollment.NewListMineHandler(uc.Enrollment.ListMine))
	// Self-enrolment is deliberately absent. Campus serves institutions that
	// decide who attends what: a published course is not an open one, and a
	// student adding themselves would enrol past whatever the institution
	// arranged. Enrolment comes from the course's teacher or its admin.
	teacherOnly.GET("/me/practiq-students", handlerProfile.NewListMyPractiqStudentsHandler(uc.Profile.ListMyPractiqStudents))
	teacherOnly.DELETE("/enrollments/:id", handlerEnrollment.NewDeleteHandler(uc.Enrollment.Delete))

	// Sections.
	teacherOnly.POST("/courses/:id/sections", handlerSection.NewCreateHandler(uc.CourseSection.Create))
	api.GET("/courses/:id/sections", handlerSection.NewListHandler(uc.CourseSection.List))
	teacherOnly.PUT("/sections/:id", handlerSection.NewUpdateHandler(uc.CourseSection.Update))
	teacherOnly.DELETE("/sections/:id", handlerSection.NewDeleteHandler(uc.CourseSection.Delete))

	// Groups — a course-scoped cohort ("Comisión A") for organizing and
	// filtering enrolled students. Reading is open to the plain group like
	// sections; managing membership stays teacher-only.
	api.GET("/courses/:id/groups", handlerGroup.List(uc.CourseGroup.Manage))
	teacherOnly.POST("/courses/:id/groups", handlerGroup.Create(uc.CourseGroup.Manage))
	teacherOnly.DELETE("/groups/:id", handlerGroup.Delete(uc.CourseGroup.Manage))
	teacherOnly.POST("/groups/:id/members", handlerGroup.AddMember(uc.CourseGroup.Manage))
	teacherOnly.DELETE("/groups/:id/members/:userId", handlerGroup.RemoveMember(uc.CourseGroup.Manage))

	// Materials — a file is uploaded first (POST /uploads returns a bucket
	// URL), then referenced here; an external link skips the upload step.
	// Reading is open to enrolled students, so it hangs off the plain group.
	api.POST("/uploads", handlerUpload.NewHandler(uc.Upload.Upload))
	teacherOnly.POST("/courses/:id/materials", handlerMaterial.NewCreateHandler(uc.CourseMaterial.Create))
	api.GET("/courses/:id/materials", handlerMaterial.NewListHandler(uc.CourseMaterial.List))
	teacherOnly.DELETE("/materials/:id", handlerMaterial.NewDeleteHandler(uc.CourseMaterial.Delete))

	// Assignments.
	teacherOnly.POST("/courses/:id/assignments", handlerAssignment.NewCreateHandler(uc.Assignment.Create))
	api.GET("/courses/:id/assignments", handlerAssignment.NewListHandler(uc.Assignment.List))
	teacherOnly.PUT("/assignments/:id", handlerAssignment.NewUpdateHandler(uc.Assignment.Update))
	teacherOnly.DELETE("/assignments/:id", handlerAssignment.NewDeleteHandler(uc.Assignment.Delete))
	api.GET("/assignments/:id/rubric", handlerRubric.List(uc.Rubric.Manage))
	teacherOnly.PUT("/assignments/:id/rubric", handlerRubric.Replace(uc.Rubric.Manage))

	// Submissions — creating one is gated by enrollment (checked inside the
	// usecase), not by role, so it hangs off the plain api group.
	api.POST("/assignments/:id/submissions", handlerSubmission.NewCreateHandler(uc.Submission.Create))
	api.GET("/assignments/:id/submissions/me", handlerSubmission.NewGetMineHandler(uc.Submission.GetMine))
	teacherOnly.GET("/assignments/:id/submissions", handlerSubmission.NewListByAssignmentHandler(uc.Submission.ListByAssignment))
	teacherOnly.PUT("/submissions/:id/grade", handlerSubmission.NewGradeHandler(uc.Submission.Grade))

	// Quizzes — auto-graded (multiple choice, true/false, fill-in-the-blanks),
	// deliberately never AI: every answer is decided by exact comparison so a
	// result is reproducible and available the instant a student submits.
	teacherOnly.POST("/courses/:id/quizzes", handlerQuiz.NewCreateHandler(uc.Quiz.Create))
	api.GET("/courses/:id/quizzes", handlerQuiz.NewListHandler(uc.Quiz.List))
	api.GET("/quizzes/:id", handlerQuiz.NewGetHandler(uc.Quiz.Get))
	teacherOnly.PUT("/quizzes/:id", handlerQuiz.NewUpdateHandler(uc.Quiz.Update))
	teacherOnly.DELETE("/quizzes/:id", handlerQuiz.NewDeleteHandler(uc.Quiz.Delete))
	teacherOnly.GET("/quizzes/:id/questions", handlerQuiz.NewListQuestionsHandler(uc.Quiz.Questions))
	teacherOnly.PUT("/quizzes/:id/questions", handlerQuiz.NewReplaceQuestionsHandler(uc.Quiz.Questions))

	// Quiz attempts — starting one is gated by enrollment (checked inside the
	// usecase), not by role, so it hangs off the plain api group.
	api.POST("/quizzes/:id/attempts", handlerQuizAttempt.NewStartHandler(uc.QuizAttempt.Start))
	api.GET("/quizzes/:id/attempts/mine", handlerQuizAttempt.NewListMineHandler(uc.QuizAttempt.ListMine))
	teacherOnly.GET("/quizzes/:id/attempts", handlerQuizAttempt.NewListByQuizHandler(uc.QuizAttempt.ListByQuiz))
	api.POST("/attempts/:id/submit", handlerQuizAttempt.NewSubmitHandler(uc.QuizAttempt.Submit))
	api.GET("/attempts/:id", handlerQuizAttempt.NewGetHandler(uc.QuizAttempt.Get))

	// Forum — access (owner, superadmin, or enrolled) is checked inside the
	// usecases, not by role, since students post too.
	// Only course teachers can open new discussion topics. Enrolled students
	// can read topics and reply through the plain authenticated routes below.
	teacherOnly.POST("/courses/:id/forum-threads", handlerThread.NewCreateHandler(uc.ForumThread.Create))
	teacherOnly.PUT("/forum-threads/:id", handlerThread.NewUpdateHandler(uc.ForumThread.Update))
	api.GET("/courses/:id/forum-threads", handlerThread.NewListHandler(uc.ForumThread.List))
	api.POST("/forum-threads/:id/posts", handlerPost.NewCreateHandler(uc.ForumPost.Create))
	api.GET("/forum-threads/:id/posts", handlerPost.NewListHandler(uc.ForumPost.List))

	// Calendar — teachers create course events and select enrolled attendees.
	teacherOnly.POST("/calendar/events", handlerCalendar.NewCreateHandler(uc.CalendarEvent.Create))
	teacherOnly.PUT("/calendar/events/:id", handlerCalendar.NewUpdateHandler(uc.CalendarEvent.Update))
	teacherOnly.DELETE("/calendar/events/:id", handlerCalendar.NewDeleteHandler(uc.CalendarEvent.Delete))
	api.GET("/me/calendar", handlerCalendar.NewListMineHandler(uc.CalendarEvent.ListMine))

	// Messaging — a 1:1 DM by the recipient's email, gated by "do these two
	// people share a course" (checked inside the usecase, not by role, so
	// student-student and student-teacher both work the same way).
	api.POST("/messages", handlerMessage.NewSendHandler(uc.Message.Send))
	api.GET("/courses/:id/messageable-users", handlerMessage.NewSearchRecipientsHandler(uc.Message.SearchRecipients))
	teacherOnly.POST("/courses/:id/messages/broadcast", handlerMessage.NewBroadcastHandler(uc.Message.Broadcast))
	api.GET("/conversations", handlerMessage.NewListConversationsHandler(uc.Message.ListConversations))
	api.GET("/conversations/:id/messages", handlerMessage.NewListMessagesHandler(uc.Message.ListMessages))
	api.POST("/conversations/:id/messages", handlerMessage.NewReplyHandler(uc.Message.Reply))
	api.POST("/conversations/:id/read", handlerMessage.NewMarkReadHandler(uc.Message.MarkRead))
}
