package tenantcontext

import (
	"context"
	"strings"
	"testing"
)

func withTenant(id string) context.Context {
	return context.WithValue(context.Background(), Key, id)
}

// The invariant every other test rests on. A statement built without a tenant
// must not be runnable at all — not run unfiltered, not run matching nothing.
func TestQueryRefusesToBuildWithoutATenant(t *testing.T) {
	_, _, err := NewQuery(context.Background()).Own("c.tenant_id").SQL("c.id", "courses c")
	if err != ErrNoTenant {
		t.Fatalf("err = %v, want ErrNoTenant", err)
	}
}

func TestOwnConstrainsTheTableItself(t *testing.T) {
	sql, args, err := NewQuery(withTenant("t1"), "course-1").
		Where("c.id = ?", "course-1").
		Own("c.tenant_id").
		SQL("c.id, c.title", "courses c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "c.tenant_id = $") {
		t.Fatalf("statement does not constrain the tenant: %s", sql)
	}
	if args[len(args)-1] != "t1" {
		t.Fatalf("tenant not bound, args = %v", args)
	}
}

// A child row is reached through the ancestor that carries the tenant, and the
// constraint travels with the join rather than being a separate check.
func TestThroughConstrainsViaTheAncestor(t *testing.T) {
	sql, _, err := NewQuery(withTenant("t1")).
		Through("courses", "c", "a.course_id").
		Where("a.id = ?", "assignment-1").
		SQL("a.id", "assignments a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "JOIN courses c ON c.id = a.course_id AND c.tenant_id = $") {
		t.Fatalf("ancestor join is not tenant-constrained: %s", sql)
	}
}

// Two hops: a submission belongs to a tenant only through its assignment's
// course. Neither table in between carries tenant_id.
func TestChainWalksToTheTenantOwningTable(t *testing.T) {
	sql, _, err := NewQuery(withTenant("t1")).
		Join("assignments", "a", "s.assignment_id").
		Through("courses", "c", "a.course_id").
		Where("s.id = ?", "submission-1").
		SQL("s.id", "submissions s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "JOIN assignments a ON a.id = s.assignment_id") {
		t.Fatalf("missing intermediate join: %s", sql)
	}
	if !strings.Contains(sql, "c.tenant_id = $") {
		t.Fatalf("chain does not end at the tenant: %s", sql)
	}
}

// Placeholders have to keep counting across the caller's arguments and the
// tenant's, or the wrong value binds to the wrong column.
func TestPlaceholdersDoNotCollide(t *testing.T) {
	sql, args, err := NewQuery(withTenant("t1")).
		Where("a.id = ?", "a1").
		Through("courses", "c", "a.course_id").
		Where("a.status = ?", "open").
		SQL("a.id", "assignments a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, want := range []any{"a1", "t1", "open"} {
		if args[i] != want {
			t.Fatalf("args[%d] = %v, want %v — binding order drifted: %s", i, args[i], want, sql)
		}
	}
	for _, placeholder := range []string{"$1", "$2", "$3"} {
		if !strings.Contains(sql, placeholder) {
			t.Fatalf("missing %s in %s", placeholder, sql)
		}
	}
}
