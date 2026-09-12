// Package campusaccess decides who may manage what inside one institution.
//
// It exists so the rule is written once. The same question — "may this person
// change this course?" — was answered in a dozen use cases by comparing
// owner_id, which meant an institution's own administrator could not touch the
// courses of the teachers they administer, and that adding a role later would
// have meant finding every copy.
package campusaccess

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// CanManageCourse reports whether the caller may change a course.
//
// Three people can: the teacher who owns it, the administrator of the
// institution it belongs to, and a platform superadmin. The tenant is already
// settled by the time this runs — the course was read through a scoped query —
// so this only decides rank within that institution, never which one.
func CanManageCourse(ctx context.Context, ownerID, requesterID string, isSuperAdmin bool) bool {
	if isSuperAdmin {
		return true
	}
	if ownerID != "" && ownerID == requesterID {
		return true
	}
	// An institution's administrator runs its courses whoever teaches them.
	// That is what being its administrator means, and it is scoped: an admin
	// of one institution is an ordinary member of the rest.
	return tenantcontext.IsAdmin(ctx)
}

// CanManageInstitution reports whether the caller may act on the institution
// itself — its people and its settings — rather than on one course inside it.
func CanManageInstitution(ctx context.Context, isSuperAdmin bool) bool {
	return isSuperAdmin || tenantcontext.IsAdmin(ctx)
}

// IsStudent reports whether the caller is only a student here.
//
// Used for the reads that must stay narrow: their own enrolments, their own
// work. Written as "is a student" rather than "is not a teacher" so a caller
// with no role at all — a request that never resolved one — is not mistaken
// for staff.
func IsStudent(ctx context.Context) bool {
	return tenantcontext.Role(ctx) == tenantcontext.RoleStudent
}
