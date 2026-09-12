package practiqapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type (
	ProfileInfo struct {
		ID          string
		Name        string
		ProfileType string
	}

	StudentInfo struct {
		ID    string
		Name  string
		Email string
	}

	SchoolInfo struct {
		ID      string
		Name    string
		Kind    string
		Billing string
		Status  string
		Role    string
	}

	SchoolMemberInfo struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Active bool   `json:"active"`
	}

	SubjectInfo struct {
		ID          string
		Name        string
		Description string
		CreatedBy   string
	}

	// Client talks to practiq-be's own profile endpoint, so Campus can tell
	// whether a shared identity is already known there as a student or a
	// teacher before deciding how to provision it locally.
	Client interface {
		// GetProfile returns nil, nil when practiq-be has no profile for
		// that id — not found is not an error here.
		GetProfile(ctx context.Context, bearerToken, id string) (*ProfileInfo, error)
		// ListMyStudents forwards the caller's own bearer token to
		// practiq-be's teacher-student-assignments — it always returns that
		// teacher's own list, nothing else.
		ListMyStudents(ctx context.Context, bearerToken string) ([]StudentInfo, error)
		// ListSubjects returns practiq-be's full subject catalog (it has no
		// per-teacher filter of its own) — callers filter by CreatedBy
		// themselves when they only want "my" subjects.
		ListSubjects(ctx context.Context, bearerToken string) ([]SubjectInfo, error)
		ListMySchools(ctx context.Context, bearerToken string) ([]SchoolInfo, error)
		// ListAllSchools is practiq-be's platform-wide listing, which it only
		// serves to a superadmin. Campus uses it to confirm a school exists
		// and qualifies before enabling it, rather than trusting the id it was
		// handed. Read-only: practiq-be stays the owner of school records.
		ListAllSchools(ctx context.Context, bearerToken string) ([]SchoolInfo, error)
		ListSchoolMembers(ctx context.Context, bearerToken, schoolID string) ([]SchoolMemberInfo, error)
		AddSchoolMember(ctx context.Context, bearerToken, schoolID, userID, role string) error
		RemoveSchoolMember(ctx context.Context, bearerToken, schoolID, userID string) error
	}

	client struct {
		baseURL string
		http    *http.Client
	}

	UpstreamError struct {
		Status  int
		Code    string
		Message string
	}
)

func (e *UpstreamError) Error() string {
	return fmt.Sprintf("practiq-be refused (status %d): %s", e.Status, e.Message)
}

func upstreamError(resp *http.Response, fallback string) error {
	var parsed struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	if parsed.Message == "" {
		parsed.Message = fallback
	}
	return &UpstreamError{Status: resp.StatusCode, Code: parsed.Code, Message: parsed.Message}
}

func NewClient(baseURL string) Client {
	return &client{baseURL: baseURL, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *client) ListMySchools(ctx context.Context, bearerToken string) ([]SchoolInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/schools/mine", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("practiq-be school scope failed (status %d)", resp.StatusCode)
	}
	var parsed struct {
		Data []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Kind    string `json:"kind"`
			Billing string `json:"billing"`
			Status  string `json:"status"`
			Role    string `json:"role"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	result := make([]SchoolInfo, 0, len(parsed.Data))
	for _, s := range parsed.Data {
		result = append(result, SchoolInfo{ID: s.ID, Name: s.Name, Kind: s.Kind, Billing: s.Billing, Status: s.Status, Role: s.Role})
	}
	return result, nil
}

// ListAllSchools reads every school practiq-be knows. Authorisation is
// practiq-be's: it refuses this to anyone who is not a platform superadmin, so
// Campus does not get a second, weaker copy of that rule.
func (c *client) ListAllSchools(ctx context.Context, bearerToken string) ([]SchoolInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/schools", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("practiq-be school listing failed (status %d)", resp.StatusCode)
	}
	var parsed struct {
		Data []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Kind    string `json:"kind"`
			Billing string `json:"billing"`
			Status  string `json:"status"`
			Role    string `json:"role"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	result := make([]SchoolInfo, 0, len(parsed.Data))
	for _, s := range parsed.Data {
		result = append(result, SchoolInfo{ID: s.ID, Name: s.Name, Kind: s.Kind, Billing: s.Billing, Status: s.Status, Role: s.Role})
	}
	return result, nil
}

func (c *client) GetProfile(ctx context.Context, bearerToken, id string) (*ProfileInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/profile/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("practiq-be profile lookup failed (status %d): %v", resp.StatusCode, errBody)
	}

	var parsed struct {
		Data struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ProfileType string `json:"profile_type"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &ProfileInfo{
		ID:          parsed.Data.ID,
		Name:        parsed.Data.Name,
		ProfileType: parsed.Data.ProfileType,
	}, nil
}

func (c *client) ListMyStudents(ctx context.Context, bearerToken string) ([]StudentInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/teachers/me/students", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("practiq-be my-students lookup failed (status %d): %v", resp.StatusCode, errBody)
	}

	var parsed struct {
		Data []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	students := make([]StudentInfo, 0, len(parsed.Data))
	for _, s := range parsed.Data {
		students = append(students, StudentInfo{ID: s.ID, Name: s.Name, Email: s.Email})
	}
	return students, nil
}

func (c *client) ListSubjects(ctx context.Context, bearerToken string) ([]SubjectInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/subjects", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("practiq-be subjects lookup failed (status %d): %v", resp.StatusCode, errBody)
	}

	var parsed struct {
		Data []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			CreatedBy   string `json:"created_by"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	subjects := make([]SubjectInfo, 0, len(parsed.Data))
	for _, s := range parsed.Data {
		subjects = append(subjects, SubjectInfo{ID: s.ID, Name: s.Name, Description: s.Description, CreatedBy: s.CreatedBy})
	}
	return subjects, nil
}

// ListSchoolMembers keeps School as source of truth for Campus people. Campus
// owns course data, not institution memberships or their roles.
func (c *client) ListSchoolMembers(ctx context.Context, bearerToken, schoolID string) ([]SchoolMemberInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/schools/"+url.PathEscape(schoolID)+"/members", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, upstreamError(resp, "could not load institution members")
	}
	var parsed struct {
		Data []struct {
			UserID string `json:"user_id"`
			Name   string `json:"name"`
			Email  string `json:"email"`
			Role   string `json:"role"`
			Active bool   `json:"active"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	result := make([]SchoolMemberInfo, 0, len(parsed.Data))
	for _, m := range parsed.Data {
		result = append(result, SchoolMemberInfo{UserID: m.UserID, Name: m.Name, Email: m.Email, Role: m.Role, Active: m.Active})
	}
	return result, nil
}

func (c *client) AddSchoolMember(ctx context.Context, bearerToken, schoolID, userID, role string) error {
	body, err := json.Marshal(map[string]string{"user_id": userID, "role": role})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/schools/"+url.PathEscape(schoolID)+"/members", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return upstreamError(resp, "could not add member to institution")
	}
	return nil
}

func (c *client) RemoveSchoolMember(ctx context.Context, bearerToken, schoolID, userID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/api/schools/"+url.PathEscape(schoolID)+"/members/"+url.PathEscape(userID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bearerToken)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return upstreamError(resp, "could not remove member from institution")
	}
	return nil
}
