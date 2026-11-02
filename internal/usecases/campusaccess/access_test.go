package campusaccess

import (
	"context"
	"testing"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func asRole(role string) context.Context {
	return context.WithValue(context.Background(), tenantcontext.RoleKey, role)
}

func TestCanManageCourse(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		ownerID      string
		requesterID  string
		isSuperAdmin bool
		want         bool
	}{
		{
			name:        "the teacher who owns the course may manage it",
			ctx:         asRole(tenantcontext.RoleTeacher),
			ownerID:     "teacher-1",
			requesterID: "teacher-1",
			want:        true,
		},
		{

			name:        "the institution's admin may manage a course they do not own",
			ctx:         asRole(tenantcontext.RoleAdmin),
			ownerID:     "teacher-1",
			requesterID: "admin-1",
			want:        true,
		},
		{
			name:        "another teacher in the same institution may not",
			ctx:         asRole(tenantcontext.RoleTeacher),
			ownerID:     "teacher-1",
			requesterID: "teacher-2",
		},
		{
			name:        "a student may not",
			ctx:         asRole(tenantcontext.RoleStudent),
			ownerID:     "teacher-1",
			requesterID: "student-1",
		},
		{
			name:         "a platform superadmin may",
			ctx:          asRole(tenantcontext.RoleStudent),
			ownerID:      "teacher-1",
			requesterID:  "root",
			isSuperAdmin: true,
			want:         true,
		},
		{

			name:        "no role and no owner grants nothing",
			ctx:         context.Background(),
			ownerID:     "",
			requesterID: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CanManageCourse(tc.ctx, tc.ownerID, tc.requesterID, tc.isSuperAdmin)
			if got != tc.want {
				t.Fatalf("CanManageCourse = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAdminRightsDoNotTravelBetweenInstitutions(t *testing.T) {

	if CanManageCourse(context.Background(), "teacher-1", "admin-of-another-school", false) {
		t.Fatal("a role held in another institution must not grant anything here")
	}
	if CanManageInstitution(context.Background(), false) {
		t.Fatal("administering one institution is not administering every institution")
	}
}
