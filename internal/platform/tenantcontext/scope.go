package tenantcontext

import (
	"context"
	"errors"
)

// ErrNoTenant is returned when a query runs without a tenant in context.
//
// It is an error and not an empty filter on purpose. A missing tenant means
// the request never passed RequireTenant, and answering it with unfiltered
// rows would turn a routing mistake into a cross-institution leak. Failing
// here costs a 500 on a path that should not exist; failing open costs
// somebody else's data.
var ErrNoTenant = errors.New("tenant is not in context")

// Scope is the tenant filter for a query, closed by default.
//
// Every repository composes its WHERE through this rather than reading the
// tenant itself, so there is one definition of "belongs to this tenant" and
// one place where its absence is refused.
type Scope struct {
	// Condition is SQL that constrains rows to the tenant, already written
	// against the column handed to Require.
	Condition string
	// Arg is the tenant id the condition binds.
	Arg any
}

// Require builds the filter for a tenant-owned column, refusing a request
// with no tenant.
//
// column is qualified by the caller ("c.tenant_id", "courses.tenant_id") so
// the same helper serves a plain SELECT and one joined to its parent.
// placeholder is the positional parameter the caller has room for.
func Require(ctx context.Context, column, placeholder string) (Scope, error) {
	id := ID(ctx)
	if id == "" {
		return Scope{}, ErrNoTenant
	}
	return Scope{Condition: column + " = " + placeholder, Arg: id}, nil
}
