package domain

import "time"

type Notification struct {
	ID, UserID, Type, Title, Body, Data string
	ReadAt                              *time.Time
	CreatedAt                           time.Time
}
