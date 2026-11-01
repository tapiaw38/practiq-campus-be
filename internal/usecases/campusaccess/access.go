package campusaccess

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func CanManageCourse(ctx context.Context, ownerID, requesterID string, isSuperAdmin bool) bool {
	if isSuperAdmin {
		return true
	}
	if ownerID != "" && ownerID == requesterID {
		return true
	}

	return tenantcontext.IsAdmin(ctx)
}

func CanManageInstitution(ctx context.Context, isSuperAdmin bool) bool {
	return isSuperAdmin || tenantcontext.IsAdmin(ctx)
}

func IsStudent(ctx context.Context) bool {
	return tenantcontext.Role(ctx) == tenantcontext.RoleStudent
}
