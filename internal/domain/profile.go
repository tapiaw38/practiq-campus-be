package domain

import "time"

// Profile mirrors practiq-be's per-app local profile: ID is not a generated
// key, it IS the user_id auth-api-be issues in the JWT, so a lookup never
// needs a join — same shape as practiq-be's user_profiles table.
//
// Name/email are deliberately NOT stored here — auth-api-be is the single
// source of truth for identity. Callers that need display name/email fetch
// it from auth-api-be (see internal/platform/identity) rather than trusting
// a local cache that can drift.
type Profile struct {
	ID          string
	ProfileType string // "student" | "teacher"
	AvatarURL   string
	Bio         string
	IsBlocked   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
