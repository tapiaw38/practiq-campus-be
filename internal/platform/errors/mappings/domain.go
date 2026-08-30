package mappings

import "net/http"

var (
	ProfileSyncError = ErrorDetails{
		InternalCode: "profile:sync-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to sync profile",
	}
	ProfileGetError = ErrorDetails{
		InternalCode: "profile:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get profile",
	}
	ProfileNotFoundError = ErrorDetails{
		InternalCode: "profile:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "profile not found",
	}
	ProfileEmailTakenError = ErrorDetails{
		InternalCode: "profile:email-taken",
		StatusCode:   http.StatusConflict,
		Message:      "a profile with this email already exists",
	}
	ProfileProvisionAuthError = ErrorDetails{
		InternalCode: "profile:provision-auth-error",
		StatusCode:   http.StatusBadGateway,
		Message:      "failed to create the shared account",
	}
	PractiqLookupError = ErrorDetails{
		InternalCode: "profile:practiq-lookup-error",
		StatusCode:   http.StatusBadGateway,
		Message:      "failed to reach practiq",
	}
	ProfileMissingFieldsError = ErrorDetails{
		InternalCode: "profile:missing-fields",
		StatusCode:   http.StatusUnprocessableEntity,
		Message:      "this email has no account yet — first_name, last_name and password are required to create one",
	}
	PreferenceGetError = ErrorDetails{
		InternalCode: "preference:get-error", StatusCode: http.StatusInternalServerError, Message: "failed to get preferences",
	}
	PreferenceUpdateError = ErrorDetails{
		InternalCode: "preference:update-error", StatusCode: http.StatusInternalServerError, Message: "failed to update preferences",
	}

	CourseCreateError = ErrorDetails{
		InternalCode: "course:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create course",
	}
	CourseListError = ErrorDetails{
		InternalCode: "course:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list courses",
	}
	CourseGetError = ErrorDetails{
		InternalCode: "course:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get course",
	}
	CourseUpdateError = ErrorDetails{
		InternalCode: "course:update-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to update course",
	}
	CourseNotFoundError = ErrorDetails{
		InternalCode: "course:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "course not found",
	}
	CourseSlugTakenError = ErrorDetails{
		InternalCode: "course:slug-taken",
		StatusCode:   http.StatusConflict,
		Message:      "a course with this slug already exists",
	}

	EnrollmentCreateError = ErrorDetails{
		InternalCode: "enrollment:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create enrollment",
	}
	EnrollmentListError = ErrorDetails{
		InternalCode: "enrollment:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list enrollments",
	}
	EnrollmentDeleteError = ErrorDetails{
		InternalCode: "enrollment:delete-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to delete enrollment",
	}
	EnrollmentGetError = ErrorDetails{
		InternalCode: "enrollment:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get enrollment",
	}
	EnrollmentNotFoundError = ErrorDetails{
		InternalCode: "enrollment:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "enrollment not found",
	}
	EnrollmentAlreadyExistsError = ErrorDetails{
		InternalCode: "enrollment:already-exists",
		StatusCode:   http.StatusConflict,
		Message:      "student is already enrolled in this course",
	}
	SectionCreateError = ErrorDetails{
		InternalCode: "section:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create section",
	}
	SectionGetError = ErrorDetails{
		InternalCode: "section:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get section",
	}
	SectionListError = ErrorDetails{
		InternalCode: "section:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list sections",
	}

	AssignmentCreateError = ErrorDetails{
		InternalCode: "assignment:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create assignment",
	}
	AssignmentGetError = ErrorDetails{
		InternalCode: "assignment:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get assignment",
	}
	AssignmentListError = ErrorDetails{
		InternalCode: "assignment:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list assignments",
	}
	AssignmentNotFoundError = ErrorDetails{
		InternalCode: "assignment:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "assignment not found",
	}

	QuizCreateError = ErrorDetails{
		InternalCode: "quiz:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create quiz",
	}
	QuizGetError = ErrorDetails{
		InternalCode: "quiz:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get quiz",
	}
	QuizListError = ErrorDetails{
		InternalCode: "quiz:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list quizzes",
	}
	QuizNotFoundError = ErrorDetails{
		InternalCode: "quiz:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "quiz not found",
	}
	QuizAttemptNotFoundError = ErrorDetails{
		InternalCode: "quiz:attempt-not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "attempt not found",
	}
	QuizAttemptAlreadySubmittedError = ErrorDetails{
		InternalCode: "quiz:attempt-already-submitted",
		StatusCode:   http.StatusBadRequest,
		Message:      "this attempt was already submitted",
	}
	QuizAttemptsExhaustedError = ErrorDetails{
		InternalCode: "quiz:attempts-exhausted",
		StatusCode:   http.StatusBadRequest,
		Message:      "no attempts left for this quiz",
	}
	QuizNotAvailableError = ErrorDetails{
		InternalCode: "quiz:not-available",
		StatusCode:   http.StatusBadRequest,
		Message:      "this quiz is not open right now",
	}
	QuizNotEnrolledError = ErrorDetails{
		InternalCode: "quiz:not-enrolled",
		StatusCode:   http.StatusForbidden,
		Message:      "only enrolled students can take this quiz",
	}

	SubmissionCreateError = ErrorDetails{
		InternalCode: "submission:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create submission",
	}
	SubmissionGetError = ErrorDetails{
		InternalCode: "submission:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get submission",
	}
	SubmissionListError = ErrorDetails{
		InternalCode: "submission:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list submissions",
	}
	SubmissionNotFoundError = ErrorDetails{
		InternalCode: "submission:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "submission not found",
	}
	SubmissionNotEnrolledError = ErrorDetails{
		InternalCode: "submission:not-enrolled",
		StatusCode:   http.StatusForbidden,
		Message:      "only enrolled students can submit",
	}
	SubmissionAlreadyExistsError = ErrorDetails{
		InternalCode: "submission:already-exists",
		StatusCode:   http.StatusConflict,
		Message:      "you already submitted this assignment",
	}

	ThreadCreateError = ErrorDetails{
		InternalCode: "thread:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create thread",
	}
	ThreadGetError = ErrorDetails{
		InternalCode: "thread:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get thread",
	}
	ThreadListError = ErrorDetails{
		InternalCode: "thread:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list threads",
	}
	ThreadUpdateError = ErrorDetails{
		InternalCode: "thread:update-error", StatusCode: http.StatusInternalServerError, Message: "failed to update thread",
	}
	ThreadNotFoundError = ErrorDetails{
		InternalCode: "thread:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "thread not found",
	}

	PostCreateError = ErrorDetails{
		InternalCode: "post:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create post",
	}
	PostListError = ErrorDetails{
		InternalCode: "post:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list posts",
	}

	CalendarEventCreateError = ErrorDetails{
		InternalCode: "calendar:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create event",
	}
	CalendarEventListError = ErrorDetails{
		InternalCode: "calendar:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list calendar",
	}
	CalendarEventGetError = ErrorDetails{
		InternalCode: "calendar:get-error", StatusCode: http.StatusInternalServerError, Message: "failed to get event",
	}
	CalendarEventUpdateError = ErrorDetails{
		InternalCode: "calendar:update-error", StatusCode: http.StatusInternalServerError, Message: "failed to update event",
	}
	CalendarEventDeleteError = ErrorDetails{
		InternalCode: "calendar:delete-error", StatusCode: http.StatusInternalServerError, Message: "failed to delete event",
	}

	ConversationCreateError = ErrorDetails{
		InternalCode: "conversation:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create conversation",
	}
	ConversationListError = ErrorDetails{
		InternalCode: "conversation:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list conversations",
	}
	ConversationNotFoundError = ErrorDetails{
		InternalCode: "conversation:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "conversation not found",
	}
	MessageCreateError = ErrorDetails{
		InternalCode: "message:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to send message",
	}
	MessageListError = ErrorDetails{
		InternalCode: "message:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list messages",
	}
	MessageRecipientNotFoundError = ErrorDetails{
		InternalCode: "message:recipient-not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "no Campus account exists for that email yet",
	}
	MessageNoSharedCourseError = ErrorDetails{
		InternalCode: "message:no-shared-course",
		StatusCode:   http.StatusForbidden,
		Message:      "you can only message someone you share a course with",
	}

	UploadError = ErrorDetails{
		InternalCode: "upload:error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to upload the file",
	}
	UploadTooLargeError = ErrorDetails{
		InternalCode: "upload:too-large",
		StatusCode:   http.StatusRequestEntityTooLarge,
		Message:      "the file is too large",
	}
	UploadUnsupportedTypeError = ErrorDetails{
		InternalCode: "upload:unsupported-type",
		StatusCode:   http.StatusUnsupportedMediaType,
		Message:      "unsupported file type",
	}
	UploadNotConfiguredError = ErrorDetails{
		InternalCode: "upload:not-configured",
		StatusCode:   http.StatusServiceUnavailable,
		Message:      "file storage is not configured",
	}

	MaterialCreateError = ErrorDetails{
		InternalCode: "material:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create material",
	}
	MaterialGetError = ErrorDetails{
		InternalCode: "material:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get material",
	}
	MaterialListError = ErrorDetails{
		InternalCode: "material:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list materials",
	}
	MaterialDeleteError = ErrorDetails{
		InternalCode: "material:delete-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to delete material",
	}
	MaterialNotFoundError = ErrorDetails{
		InternalCode: "material:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "material not found",
	}
	MaterialFileNotOwnedError = ErrorDetails{
		InternalCode: "material:file-not-owned",
		StatusCode:   http.StatusForbidden,
		Message:      "the file does not belong to you",
	}

	EnrollmentStudentNotFoundError = ErrorDetails{
		InternalCode: "enrollment:student-not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "no student with this email exists in Campus yet",
	}
)
