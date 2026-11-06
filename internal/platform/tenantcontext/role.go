package tenantcontext

import "context"

const RoleKey = "campusTenantRole"

const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
)

func Role(ctx context.Context) string {
	role, _ := ctx.Value(RoleKey).(string)
	return role
}

func IsAdmin(ctx context.Context) bool {
	return Role(ctx) == RoleAdmin
}
