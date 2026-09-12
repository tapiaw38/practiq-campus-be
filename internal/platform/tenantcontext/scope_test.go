package tenantcontext

import (
	"context"
	"testing"
)

// The invariant the whole isolation rests on: no tenant is not "every tenant".
func TestRequireRefusesAMissingTenant(t *testing.T) {
	if _, err := Require(context.Background(), "c.tenant_id", "$1"); err != ErrNoTenant {
		t.Fatalf("err = %v, want ErrNoTenant — an unscoped query must not run", err)
	}
}

func TestRequireBindsTheTenantInContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), Key, "tenant-a")

	scope, err := Require(ctx, "c.tenant_id", "$2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scope.Condition != "c.tenant_id = $2" {
		t.Fatalf("condition = %q", scope.Condition)
	}
	if scope.Arg != "tenant-a" {
		t.Fatalf("arg = %v, want tenant-a", scope.Arg)
	}
}

// An empty string in context is the same as none: it would render
// "tenant_id = ''" and match nothing, but silently, which hides the bug.
func TestRequireTreatsAnEmptyTenantAsMissing(t *testing.T) {
	ctx := context.WithValue(context.Background(), Key, "")

	if _, err := Require(ctx, "c.tenant_id", "$1"); err != ErrNoTenant {
		t.Fatalf("err = %v, want ErrNoTenant", err)
	}
}
