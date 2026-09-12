package tenantcontext

import "context"

// RoleKey holds the caller's role inside the selected institution, resolved
// from practiq-be's school_members on every request.
//
// Campus keeps no membership table of its own. The role is read at the edge,
// carried for the length of the request, and never persisted here: somebody
// removed from an institution loses access on their next request rather than
// whenever a local copy is refreshed.
const RoleKey = "campusTenantRole"

// The roles an institution assigns. They are Practiq's, not Campus's: Campus
// reads them and applies them, and the authority to grant them stays with
// practiq-be.
const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
)

// Role reads the caller's role in the tenant. Empty when the request never
// passed RequireTenant, which is treated as no privileges anywhere.
func Role(ctx context.Context) string {
	role, _ := ctx.Value(RoleKey).(string)
	return role
}

// IsAdmin reports whether the caller administers the selected institution.
//
// This is the institution's own administrator, not a platform superadmin: it
// says "runs this school", and it is scoped to the tenant in context, so an
// admin of one institution is an ordinary member of every other.
func IsAdmin(ctx context.Context) bool {
	return Role(ctx) == RoleAdmin
}
