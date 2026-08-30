package course

import (
	"strings"
	"testing"
)

// The JOIN to enrollments (EnrolledUserID filter) brings in a second table
// that also has id/status columns. Bare "id"/"status" in the SELECT list
// parses fine but fails at query time with "column reference is ambiguous" —
// this pins every column here to carry the c. qualifier.
func TestEnrolledCourseListQueryQualifiesColumns(t *testing.T) {
	for _, col := range strings.Split(selectQualifiedCourseColumns, ",") {
		col = strings.TrimSpace(col)
		if col == "" {
			continue
		}
		if !strings.HasPrefix(col, "c.") {
			t.Fatalf("column %q is not qualified with the c. alias", col)
		}
	}
}
