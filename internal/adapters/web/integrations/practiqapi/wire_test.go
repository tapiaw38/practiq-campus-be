package practiqapi

import (
	"encoding/json"
	"testing"
)

func TestWireNames(t *testing.T) {
	cases := map[string]struct {
		value    any
		expected string
	}{
		"school member": {
			value:    SchoolMemberInfo{UserID: "u1", Name: "Nymia Tapia", Email: "n@example.com", Role: "student", Active: true},
			expected: `{"user_id":"u1","name":"Nymia Tapia","email":"n@example.com","role":"student","active":true}`,
		},
		"profile": {
			value:    ProfileInfo{ID: "u1", Name: "Nymia Tapia", ProfileType: "student"},
			expected: `{"id":"u1","name":"Nymia Tapia","profile_type":"student"}`,
		},
		"student": {
			value:    StudentInfo{ID: "u1", Name: "Nymia Tapia", Email: "n@example.com"},
			expected: `{"id":"u1","name":"Nymia Tapia","email":"n@example.com"}`,
		},
		"school": {
			value:    SchoolInfo{ID: "s1", Name: "Instituto", Kind: "institution", Billing: "direct", Status: "active", Role: "admin"},
			expected: `{"id":"s1","name":"Instituto","kind":"institution","billing":"direct","status":"active","role":"admin"}`,
		},
		"subject": {
			value:    SubjectInfo{ID: "sub1", Name: "Matemática", Description: "d", CreatedBy: "u1"},
			expected: `{"id":"sub1","name":"Matemática","description":"d","created_by":"u1"}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if string(encoded) != tc.expected {
				t.Fatalf("wire shape drifted\n got: %s\nwant: %s", encoded, tc.expected)
			}
		})
	}
}
