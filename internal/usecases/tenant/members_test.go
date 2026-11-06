package tenant

import (
	"errors"
	"strings"
	"testing"

	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
)

func TestMapUpstreamError(t *testing.T) {
	cases := map[string]struct {
		err          error
		expectedCode int
		expectedBody string
	}{
		"upstream refusal reaches the admin": {
			err:          &practiqapi.UpstreamError{Status: 400, Code: "school:plan-full", Message: "el plan no admite más alumnos"},
			expectedCode: 400,
			expectedBody: "el plan no admite más alumnos",
		},
		"upstream failure stays a bad gateway": {
			err:          &practiqapi.UpstreamError{Status: 500, Code: "school:lookup-error", Message: "boom"},
			expectedCode: 502,
			expectedBody: "could not add member to institution",
		},
		"transport failure stays a bad gateway": {
			err:          errors.New("connection refused"),
			expectedCode: 502,
			expectedBody: "could not add member to institution",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			appErr := mapUpstreamError(tc.err, "tenant:add-member-error", "could not add member to institution")

			if appErr.StatusCode() != tc.expectedCode {
				t.Fatalf("status = %d, want %d", appErr.StatusCode(), tc.expectedCode)
			}
			if !strings.Contains(appErr.Message(), tc.expectedBody) {
				t.Fatalf("message = %s, want it to contain %q", appErr.Message(), tc.expectedBody)
			}
		})
	}
}
