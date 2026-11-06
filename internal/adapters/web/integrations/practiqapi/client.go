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
		ID          string `json:"id"`
		Name        string `json:"name"`
		ProfileType string `json:"profile_type"`
	}

	StudentInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	SchoolInfo struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Kind    string `json:"kind"`
		Billing string `json:"billing"`
		Status  string `json:"status"`
		Role    string `json:"role"`
	}

	SchoolMemberInfo struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Active bool   `json:"active"`
	}

	SubjectInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedBy   string `json:"created_by"`
	}

	Client interface {
		GetProfile(ctx context.Context, bearerToken, id string) (*ProfileInfo, error)

		ListMyStudents(ctx context.Context, bearerToken string) ([]StudentInfo, error)

		ListSubjects(ctx context.Context, bearerToken string) ([]SubjectInfo, error)
		ListMySchools(ctx context.Context, bearerToken string) ([]SchoolInfo, error)

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
