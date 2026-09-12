package domain

import "time"

const (
	TenantStatusActive    = "active"
	TenantStatusSuspended = "suspended"
	TenantStatusClosed    = "closed"
)

// Tenant enables Campus for one contract-billed Practiq institution. School
// identity, lifecycle and memberships remain owned by practiq-be.
type Tenant struct {
	ID            string
	SchoolID      string
	Status        string
	ActivatedAt   time.Time
	DeactivatedAt *time.Time
	CreatedAt     time.Time
}
