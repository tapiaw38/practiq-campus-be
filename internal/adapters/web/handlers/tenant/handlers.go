package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	tenantRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/tenant"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type data struct {
	ID       string `json:"id"`
	SchoolID string `json:"school_id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func Mine(repo tenantRepo.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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
