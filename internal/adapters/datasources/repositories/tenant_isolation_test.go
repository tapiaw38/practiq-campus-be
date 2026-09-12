package repositories_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Repositories that legitimately query without a tenant filter, and why.
//
// Everything else must scope. This list is the whole exception surface: adding
// to it is a deliberate decision, which is the point — the isolation should be
// hard to drop by accident and easy to see when it is dropped on purpose.
var unscoped = map[string]string{
	// campus_tenants is the tenant registry itself. Scoping it by tenant would
	// be circular, and it is only reachable by a platform superadmin.
	"tenant": "the tenant registry, read before a tenant is resolved",
	// campus_profiles is one identity per person across every institution they
	// belong to. A profile is not tenant data; what they may see inside one is.
	"profile": "identity is global; membership decides what it reaches",
}

// Every repository that talks to the database has to scope by tenant.
//
// This is a source check rather than a behavioural one on purpose. The failure
// it guards against is not a wrong query, it is a missing condition in a new
// one — somebody adding a method months from now and not knowing the rule.
// A test that only exercised today's methods would pass while the gap opened.
func TestEveryRepositoryScopesByTenant(t *testing.T) {
	root := "."
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("cannot read repositories: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if reason, ok := unscoped[name]; ok {
			t.Logf("skipping %s: %s", name, reason)
			continue
		}

		t.Run(name, func(t *testing.T) {
			queries, scoped := false, false
			err := filepath.Walk(filepath.Join(root, name), func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
					return err
				}
				source, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				text := string(source)
				if strings.Contains(text, "QueryContext") ||
					strings.Contains(text, "QueryRowContext") ||
					strings.Contains(text, "ExecContext") {
					queries = true
				}
				if strings.Contains(text, "tenantcontext.") {
					scoped = true
				}
				return nil
			})
			if err != nil {
				t.Fatalf("walking %s: %v", name, err)
			}

			if queries && !scoped {
				t.Fatalf("%s runs database queries but never references tenantcontext.\n"+
					"Every read and write has to be constrained to the tenant in context — "+
					"scope it, or add it to the unscoped list above with the reason.", name)
			}
		})
	}
}
