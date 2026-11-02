package course

import (
	"strings"
	"testing"
)

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
