package profile

import "github.com/tapiaw38/practiq-campus-be/internal/domain"

type ProfileData struct {
	ID          string `json:"id"`
	ProfileType string `json:"profile_type"`
	FullName    string `json:"full_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	CreatedAt   string `json:"created_at"`
}

func toProfileData(p domain.Profile) ProfileData {
	return ProfileData{
		ID:          p.ID,
		ProfileType: p.ProfileType,
		FullName:    p.FullName,
		AvatarURL:   p.AvatarURL,
		Bio:         p.Bio,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
