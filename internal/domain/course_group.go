package domain

import "time"

type CourseGroup struct {
	ID        string
	CourseID  string
	Name      string
	CreatedAt time.Time
	MemberIDs []string
}
