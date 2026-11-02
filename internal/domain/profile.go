package domain

import "time"

type Profile struct {
	ID          string
	ProfileType string
	AvatarURL   string
	Bio         string
	IsBlocked   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
