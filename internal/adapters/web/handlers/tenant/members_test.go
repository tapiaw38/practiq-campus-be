package tenant

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
)

func TestForwardMemberError(t *testing.T) {
	gin.SetMode(gin.TestMode)

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
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)

			forwardMemberError(c, tc.err, "tenant:add-member-error", "could not add member to institution")

			if recorder.Code != tc.expectedCode {
				t.Fatalf("status = %d, want %d", recorder.Code, tc.expectedCode)
			}
			if !strings.Contains(recorder.Body.String(), tc.expectedBody) {
				t.Fatalf("body = %s, want it to contain %q", recorder.Body.String(), tc.expectedBody)
			}
		})
	}
}
