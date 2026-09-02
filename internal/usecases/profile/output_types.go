package profile

import "github.com/tapiaw38/practiq-campus-be/internal/domain"

type ProfileData struct {
	ID          string `json:"id"`
	ProfileType string `json:"profile_type"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	IsBlocked   bool   `json:"is_blocked"`
	CreatedAt   string `json:"created_at"`
}

// toProfileData takes name/email pre-resolved by the caller (from
// auth-api-be, see internal/platform/identity) since domain.Profile no
// longer carries identity fields.
func toProfileData(p domain.Profile, fullName, email string) ProfileData {
	return ProfileData{
		ID:          p.ID,
		ProfileType: p.ProfileType,
		FullName:    fullName,
		Email:       email,
		AvatarURL:   p.AvatarURL,
		Bio:         p.Bio,
		IsBlocked:   p.IsBlocked,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
