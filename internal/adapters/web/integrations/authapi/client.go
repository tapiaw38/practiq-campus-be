package authapi

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
	RegisterInput struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
	}

	RegisterOutput struct {
		ID       string
		Username string
		Email    string
	}

	UserInfo struct {
		ID        string
		Username  string
		FirstName string
		LastName  string
		Email     string
		Roles     []string
	}

	Client interface {
		Register(context.Context, RegisterInput) (*RegisterOutput, error)
		// GetByEmail forwards the caller's own bearer token (auth-api-be
		// requires superadmin for this lookup) and returns nil, nil when no
		// account exists for that email — not found is not an error here.
		GetByEmail(ctx context.Context, bearerToken, email string) (*UserInfo, error)
	}

	client struct {
		baseURL string
		http    *http.Client
	}
)

func NewClient(baseURL string) Client {
	return &client{baseURL: baseURL, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *client) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	body, err := json.Marshal(map[string]string{
		"first_name": input.FirstName,
		"last_name":  input.LastName,
		"email":      input.Email,
		"password":   input.Password,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/auth/register", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("auth-api-be register failed (status %d): %v", resp.StatusCode, errBody)
	}

	var parsed struct {
		Data struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &RegisterOutput{ID: parsed.Data.ID, Username: parsed.Data.Username, Email: parsed.Data.Email}, nil
}

func (c *client) GetByEmail(ctx context.Context, bearerToken, email string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/user/by-email?email="+url.QueryEscape(email), nil)
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
		return nil, fmt.Errorf("auth-api-be user lookup failed (status %d): %v", resp.StatusCode, errBody)
	}

	var parsed struct {
		Data struct {
			ID        string `json:"id"`
			Username  string `json:"username"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
			Roles     []struct {
				Name string `json:"name"`
			} `json:"roles"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(parsed.Data.Roles))
	for _, r := range parsed.Data.Roles {
		roles = append(roles, r.Name)
	}

	return &UserInfo{
		ID:        parsed.Data.ID,
		Username:  parsed.Data.Username,
		FirstName: parsed.Data.FirstName,
		LastName:  parsed.Data.LastName,
		Email:     parsed.Data.Email,
		Roles:     roles,
	}, nil
}
