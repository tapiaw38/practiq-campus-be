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
)
