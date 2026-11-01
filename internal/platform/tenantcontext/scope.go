package tenantcontext

import (
	"context"
	"errors"
)

var ErrNoTenant = errors.New("tenant is not in context")

type Scope struct {
	Condition string

	Arg any
}

func Require(ctx context.Context, column, placeholder string) (Scope, error) {
	id := ID(ctx)
	if id == "" {
		return Scope{}, ErrNoTenant
	}
	return Scope{Condition: column + " = " + placeholder, Arg: id}, nil
}
