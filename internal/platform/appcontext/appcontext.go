package appcontext

import (
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/storage"
)

type Context struct {
	Repositories *repositories.Repositories
	Integrations *integrations.Integrations
	Storage      storage.Storage
}

type Factory func() *Context

func NewFactory(repos *repositories.Repositories, integ *integrations.Integrations, store storage.Storage) Factory {
	ctx := &Context{
		Repositories: repos,
		Integrations: integ,
		Storage:      store,
	}
	return func() *Context {
		return ctx
	}
}
