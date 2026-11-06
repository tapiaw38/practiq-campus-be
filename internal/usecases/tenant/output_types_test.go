package tenant

import (
	"testing"

	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
)

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
			name:   "a personal school never qualifies",
			school: practiqapi.SchoolInfo{Kind: "personal", Billing: "subscription", Status: "active"},
		},
		{
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
