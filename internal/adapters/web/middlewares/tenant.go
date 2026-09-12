package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/tenant"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

const tenantContextKey = "campusTenantID"

// RequireTenant resolves selected tenant against Practiq on every request.
// Header is merely a selector; membership and institution lifecycle decide
// authorization. This prevents a guessed tenant UUID becoming an access key.
func RequireTenant(repo tenant.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Campus-Tenant-ID")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "tenant:required", "message": "X-Campus-Tenant-ID is required"})
			c.Abort()
			return
		}
		value, err := repo.Get(c, id)
		if err != nil || value == nil || value.Status != domain.TenantStatusActive {
			c.JSON(http.StatusForbidden, gin.H{"code": "tenant:unavailable", "message": "Campus tenant is unavailable"})
			c.Abort()
			return
		}
		if !IsSuperAdmin(c) {
			schools, err := practiq.ListMySchools(c, c.GetHeader("Authorization"))
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"code": "tenant:scope-unavailable", "message": "could not resolve school membership"})
				c.Abort()
				return
			}
			allowed := false
			for _, school := range schools {
				if school.ID == value.SchoolID && school.Kind == "institution" && school.Billing == "direct" && school.Status == "active" {
					allowed = true
					// The role travels with the membership that granted
					// access, so what the caller may do inside the
					// institution is decided by practiq-be too — Campus
					// never assigns it.
					c.Set(tenantcontext.RoleKey, school.Role)
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{"code": "tenant:forbidden", "message": "not an active member of this institution"})
				c.Abort()
				return
			}
		}
		c.Set(tenantContextKey, value.ID)
		c.Set(tenantcontext.Key, value.ID)
		c.Next()
	}
}

func GetTenantID(c *gin.Context) string {
	value, _ := c.Get(tenantContextKey)
	id, _ := value.(string)
	return id
}
