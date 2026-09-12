package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	tenantRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/tenant"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type data struct {
	ID       string `json:"id"`
	SchoolID string `json:"school_id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

// Mine lists the Campus institutions the caller may open.
//
// For everybody it is the ones they belong to. A platform superadmin belongs
// to none — they operate the platform, they are not enrolled anywhere — and
// filtering by membership left them with an empty selector and no way into an
// institution they had just enabled. They get every active one instead.
//
// This does not widen what they see inside one. RequireTenant still pins the
// request to the institution they picked, and every query stays scoped to it:
// a superadmin can enter any institution, never two at once.
func Mine(repo tenantRepo.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if middlewares.IsSuperAdmin(c) {
			listAllForSuperAdmin(c, repo, practiq)
			return
		}

		schools, err := practiq.ListMySchools(c, c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"code": "tenant:scope-unavailable", "message": "could not resolve school membership"})
			return
		}
		result := make([]data, 0, len(schools))
		for _, school := range schools {
			if school.Kind != "institution" || school.Billing != "direct" || school.Status != "active" {
				continue
			}
			value, err := repo.GetBySchoolID(c, school.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": "tenant:list-error", "message": "could not list Campus institutions"})
				return
			}
			if value != nil && value.Status == domain.TenantStatusActive {
				result = append(result, data{ID: value.ID, SchoolID: school.ID, Name: school.Name, Role: school.Role})
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}

// listAllForSuperAdmin answers with every active Campus institution, named
// from practiq-be so the picker shows a school and not a uuid.
//
// The role is left empty: a superadmin holds none inside the institution. What
// they may do there is decided by their platform role, which campusaccess
// already honours, and inventing "admin" here would put a membership in the
// response that nothing granted.
func listAllForSuperAdmin(c *gin.Context, repo tenantRepo.Repository, practiq practiqapi.Client) {
	tenants, err := repo.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "tenant:list-error", "message": "could not list Campus institutions"})
		return
	}

	names := map[string]string{}
	if schools, err := practiq.ListAllSchools(c, c.GetHeader("Authorization")); err == nil {
		for _, school := range schools {
			names[school.ID] = school.Name
		}
	}
	// A name that could not be resolved is not a reason to refuse the list:
	// the institution is still enterable, it just shows without one.

	result := make([]data, 0, len(tenants))
	for _, value := range tenants {
		if value.Status != domain.TenantStatusActive {
			continue
		}
		result = append(result, data{ID: value.ID, SchoolID: value.SchoolID, Name: names[value.SchoolID]})
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
