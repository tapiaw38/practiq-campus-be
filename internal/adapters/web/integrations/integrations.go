package integrations

import (
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/authapi"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
)

type Integrations struct {
	AuthAPI    authapi.Client
	PractiqAPI practiqapi.Client
}

func CreateIntegrations(authAPIURL, practiqAPIURL string) *Integrations {
	return &Integrations{
		AuthAPI:    authapi.NewClient(authAPIURL),
		PractiqAPI: practiqapi.NewClient(practiqAPIURL),
	}
}
