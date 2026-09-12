package tenantcontext

import "context"

const Key = "campusTenantID"

// ID reads the validated tenant stored by the HTTP middleware. Gin exposes
// request keys through context.Value, keeping repository code transport-agnostic.
func ID(ctx context.Context) string {
	value := ctx.Value(Key)
	id, _ := value.(string)
	return id
}
