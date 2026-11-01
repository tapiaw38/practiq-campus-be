package tenantcontext

import "context"

const Key = "campusTenantID"

func ID(ctx context.Context) string {
	value := ctx.Value(Key)
	id, _ := value.(string)
	return id
}
