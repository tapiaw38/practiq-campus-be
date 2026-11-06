package repositories_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var unscoped = map[string]string{

	"tenant": "the tenant registry, read before a tenant is resolved",

	"profile": "identity is global; membership decides what it reaches",
}

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
