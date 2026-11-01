package domain

import "time"

const (
	TenantStatusActive    = "active"
	TenantStatusSuspended = "suspended"
	TenantStatusClosed    = "closed"
)

type Tenant struct {
	ID            string
	SchoolID      string
	Status        string
	ActivatedAt   time.Time
	DeactivatedAt *time.Time
	CreatedAt     time.Time
}
