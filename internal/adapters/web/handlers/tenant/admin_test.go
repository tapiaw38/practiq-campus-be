package tenant

import (
	"testing"

	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
)

// Which institutions Campus may serve. Campus is sold to institutions invoiced
// by contract; anything else is a school Practiq bills per student, and
// enabling one would hand it a product nobody agreed to sell.
func TestEligible(t *testing.T) {
	institution := practiqapi.SchoolInfo{Kind: "institution", Billing: "direct", Status: "active"}

	cases := []struct {
		name   string
		school practiqapi.SchoolInfo
		want   bool
	}{
		{
			name:   "an active institution billed by contract qualifies",
			school: institution,
			want:   true,
		},
		{
			// A teacher's own school is on a per-student plan. Campus is not
			// part of it, at any size.
			name:   "a personal school never qualifies",
			school: practiqapi.SchoolInfo{Kind: "personal", Billing: "subscription", Status: "active"},
		},
		{
			// The billing column is what says somebody signed a contract. An
			// institution on a subscription has not.
			name:   "an institution on a subscription does not qualify",
			school: practiqapi.SchoolInfo{Kind: "institution", Billing: "subscription", Status: "active"},
		},
		{
			name:   "a suspended school does not qualify",
			school: practiqapi.SchoolInfo{Kind: "institution", Billing: "direct", Status: "suspended"},
		},
		{
			name:   "a closed school does not qualify",
			school: practiqapi.SchoolInfo{Kind: "institution", Billing: "direct", Status: "closed"},
		},
		{
			// The zero value is what an unknown school id decodes to. It must
			// not pass: not finding a school is not finding a valid one.
			name:   "an unknown school does not qualify",
			school: practiqapi.SchoolInfo{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := eligible(tc.school); got != tc.want {
				t.Fatalf("eligible = %v, want %v", got, tc.want)
			}
		})
	}
}
