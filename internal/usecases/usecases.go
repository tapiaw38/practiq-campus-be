package usecases

import (
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	ucAssignment "github.com/tapiaw38/practiq-campus-be/internal/usecases/assignment"
	ucCalendar "github.com/tapiaw38/practiq-campus-be/internal/usecases/calendar_event"
	ucCourse "github.com/tapiaw38/practiq-campus-be/internal/usecases/course"
	ucMaterial "github.com/tapiaw38/practiq-campus-be/internal/usecases/course_material"
	ucSection "github.com/tapiaw38/practiq-campus-be/internal/usecases/course_section"
	ucEnrollment "github.com/tapiaw38/practiq-campus-be/internal/usecases/enrollment"
	ucPost "github.com/tapiaw38/practiq-campus-be/internal/usecases/forum_post"
	ucThread "github.com/tapiaw38/practiq-campus-be/internal/usecases/forum_thread"
	ucMessage "github.com/tapiaw38/practiq-campus-be/internal/usecases/message"
	ucPreference "github.com/tapiaw38/practiq-campus-be/internal/usecases/preference"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
	ucSubmission "github.com/tapiaw38/practiq-campus-be/internal/usecases/submission"
	ucUpload "github.com/tapiaw38/practiq-campus-be/internal/usecases/upload"
)

type ProfileUsecases struct {
	Sync                  ucProfile.SyncUsecase
	Get                   ucProfile.GetUsecase
	CreateOrSync          ucProfile.CreateOrSyncUserUsecase
	ListUsers             ucProfile.ListStudentsUsecase
	ListMyPractiqStudents ucProfile.ListMyPractiqStudentsUsecase
	Block                 ucProfile.BlockUsecase
}

type PreferenceUsecases struct {
	Get    ucPreference.GetUsecase
	Update ucPreference.UpdateUsecase
}

type CourseUsecases struct {
	Create          ucCourse.CreateUsecase
	List            ucCourse.ListUsecase
	Get             ucCourse.GetUsecase
	Update          ucCourse.UpdateUsecase
	Delete          ucCourse.DeleteUsecase
	SyncFromPractiq ucCourse.SyncFromPractiqUsecase
}

type EnrollmentUsecases struct {
	Create       ucEnrollment.CreateUsecase
	ListByCourse ucEnrollment.ListByCourseUsecase
	ListMine     ucEnrollment.ListMineUsecase
	Delete       ucEnrollment.DeleteUsecase
	Self         ucEnrollment.SelfUsecase
}

type CourseMaterialUsecases struct {
	Create ucMaterial.CreateUsecase
	List   ucMaterial.ListUsecase
	Delete ucMaterial.DeleteUsecase
}

type UploadUsecases struct {
	Upload ucUpload.Usecase
}

type CourseSectionUsecases struct {
	Create ucSection.CreateUsecase
	List   ucSection.ListUsecase
	Update ucSection.UpdateUsecase
	Delete ucSection.DeleteUsecase
}

type AssignmentUsecases struct {
	Create ucAssignment.CreateUsecase
	List   ucAssignment.ListUsecase
	Update ucAssignment.UpdateUsecase
	Delete ucAssignment.DeleteUsecase
}

type SubmissionUsecases struct {
	Create           ucSubmission.CreateUsecase
	ListByAssignment ucSubmission.ListByAssignmentUsecase
	GetMine          ucSubmission.GetMineUsecase
	Grade            ucSubmission.GradeUsecase
}

type ForumThreadUsecases struct {
	Create ucThread.CreateUsecase
	List   ucThread.ListUsecase
	Update ucThread.UpdateUsecase
}

type ForumPostUsecases struct {
	Create ucPost.CreateUsecase
	List   ucPost.ListUsecase
}

type CalendarEventUsecases struct {
	Create   ucCalendar.CreateUsecase
	ListMine ucCalendar.ListMineUsecase
	Update   ucCalendar.UpdateUsecase
	Delete   ucCalendar.DeleteUsecase
}

type MessageUsecases struct {
	Send              ucMessage.SendUsecase
	Broadcast         ucMessage.BroadcastUsecase
	ListConversations ucMessage.ListConversationsUsecase
	ListMessages      ucMessage.ListMessagesUsecase
	Reply             ucMessage.ReplyUsecase
	SearchRecipients  ucMessage.SearchRecipientsUsecase
	MarkRead          ucMessage.MarkReadUsecase
}

type Usecases struct {
	Profile        ProfileUsecases
	Preference     PreferenceUsecases
	Course         CourseUsecases
	Enrollment     EnrollmentUsecases
	CourseSection  CourseSectionUsecases
	CourseMaterial CourseMaterialUsecases
	Upload         UploadUsecases
	Assignment     AssignmentUsecases
	Submission     SubmissionUsecases
	ForumThread    ForumThreadUsecases
	ForumPost      ForumPostUsecases
	CalendarEvent  CalendarEventUsecases
	Message        MessageUsecases
}

func NewUsecases(contextFactory appcontext.Factory) *Usecases {
	return &Usecases{
		Profile: ProfileUsecases{
			Sync:                  ucProfile.NewSyncUsecase(contextFactory),
			Get:                   ucProfile.NewGetUsecase(contextFactory),
			CreateOrSync:          ucProfile.NewProvisionStudentUsecase(contextFactory),
			ListUsers:             ucProfile.NewListStudentsUsecase(contextFactory),
			ListMyPractiqStudents: ucProfile.NewListMyPractiqStudentsUsecase(contextFactory),
			Block:                 ucProfile.NewBlockUsecase(contextFactory),
		},
		Preference: PreferenceUsecases{
			Get: ucPreference.NewGetUsecase(contextFactory), Update: ucPreference.NewUpdateUsecase(contextFactory),
		},
		Course: CourseUsecases{
			Create:          ucCourse.NewCreateUsecase(contextFactory),
			List:            ucCourse.NewListUsecase(contextFactory),
			Get:             ucCourse.NewGetUsecase(contextFactory),
			Update:          ucCourse.NewUpdateUsecase(contextFactory),
			Delete:          ucCourse.NewDeleteUsecase(contextFactory),
			SyncFromPractiq: ucCourse.NewSyncFromPractiqUsecase(contextFactory),
		},
		Enrollment: EnrollmentUsecases{
			Create:       ucEnrollment.NewCreateUsecase(contextFactory),
			ListByCourse: ucEnrollment.NewListByCourseUsecase(contextFactory),
			ListMine:     ucEnrollment.NewListMineUsecase(contextFactory),
			Delete:       ucEnrollment.NewDeleteUsecase(contextFactory),
			Self:         ucEnrollment.NewSelfUsecase(contextFactory),
		},
		CourseSection: CourseSectionUsecases{
			Create: ucSection.NewCreateUsecase(contextFactory),
			List:   ucSection.NewListUsecase(contextFactory),
			Update: ucSection.NewUpdateUsecase(contextFactory), Delete: ucSection.NewDeleteUsecase(contextFactory),
		},
		CourseMaterial: CourseMaterialUsecases{
			Create: ucMaterial.NewCreateUsecase(contextFactory),
			List:   ucMaterial.NewListUsecase(contextFactory),
			Delete: ucMaterial.NewDeleteUsecase(contextFactory),
		},
		Upload: UploadUsecases{
			Upload: ucUpload.NewUsecase(contextFactory),
		},
		Assignment: AssignmentUsecases{
			Create: ucAssignment.NewCreateUsecase(contextFactory),
			List:   ucAssignment.NewListUsecase(contextFactory),
			Update: ucAssignment.NewUpdateUsecase(contextFactory), Delete: ucAssignment.NewDeleteUsecase(contextFactory),
		},
		Submission: SubmissionUsecases{
			Create:           ucSubmission.NewCreateUsecase(contextFactory),
			ListByAssignment: ucSubmission.NewListByAssignmentUsecase(contextFactory),
			GetMine:          ucSubmission.NewGetMineUsecase(contextFactory),
			Grade:            ucSubmission.NewGradeUsecase(contextFactory),
		},
		ForumThread: ForumThreadUsecases{
			Create: ucThread.NewCreateUsecase(contextFactory),
			List:   ucThread.NewListUsecase(contextFactory),
			Update: ucThread.NewUpdateUsecase(contextFactory),
		},
		ForumPost: ForumPostUsecases{
			Create: ucPost.NewCreateUsecase(contextFactory),
			List:   ucPost.NewListUsecase(contextFactory),
		},
		CalendarEvent: CalendarEventUsecases{
			Create:   ucCalendar.NewCreateUsecase(contextFactory),
			ListMine: ucCalendar.NewListMineUsecase(contextFactory),
			Update:   ucCalendar.NewUpdateUsecase(contextFactory),
			Delete:   ucCalendar.NewDeleteUsecase(contextFactory),
		},
		Message: MessageUsecases{
			Send:              ucMessage.NewSendUsecase(contextFactory),
			Broadcast:         ucMessage.NewBroadcastUsecase(contextFactory),
			ListConversations: ucMessage.NewListConversationsUsecase(contextFactory),
			ListMessages:      ucMessage.NewListMessagesUsecase(contextFactory),
			Reply:             ucMessage.NewReplyUsecase(contextFactory),
			SearchRecipients:  ucMessage.NewSearchRecipientsUsecase(contextFactory),
			MarkRead:          ucMessage.NewMarkReadUsecase(contextFactory),
		},
	}
}
