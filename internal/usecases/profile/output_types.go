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

func toProfileData(p domain.Profile) ProfileData {
	return ProfileData{
		ID:          p.ID,
		ProfileType: p.ProfileType,
		FullName:    p.FullName,
		Email:       p.Email,
		AvatarURL:   p.AvatarURL,
		Bio:         p.Bio,
		IsBlocked:   p.IsBlocked,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toProfileDataList(profiles []domain.Profile) []ProfileData {
	out := make([]ProfileData, 0, len(profiles))
	for _, p := range profiles {
		out = append(out, toProfileData(p))
	}
	return out
}
