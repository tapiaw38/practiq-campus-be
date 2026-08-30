package domain

import "time"

// Profile mirrors practiq-be's per-app local profile: ID is not a generated
// key, it IS the user_id auth-api-be issues in the JWT, so a lookup never
// needs a join — same shape as practiq-be's user_profiles table.
type Profile struct {
	ID          string
	ProfileType string // "student" | "teacher"
	FullName    string
	AvatarURL   string
	Bio         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
