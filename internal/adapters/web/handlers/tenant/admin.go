package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"

	tenantRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/tenant"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

// adminData carries the school's own details alongside the Campus record, so
// an operator sees a name rather than a pair of UUIDs. The name is read from
// practiq-be every time instead of copied here: a school renamed there must
// not go on showing its old name in Campus.
type adminData struct {
	ID       string `json:"id"`
	SchoolID string `json:"school_id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	// Eligible reports whether the school still qualifies for Campus. A tenant
	// whose school was closed or moved off contract billing keeps its record
	// but stops qualifying, and an operator needs to see that.
	Eligible bool `json:"eligible"`
}

// eligible is the single definition of which institutions Campus may serve.
//
// Campus exists for institutions invoiced by contract. A personal school, or
// one on a per-student subscription, is billed by Practiq and is not a Campus
// customer — enabling one would hand it a product nobody agreed to sell.
func eligible(school practiqapi.SchoolInfo) bool {
	return school.Kind == "institution" &&
		school.Billing == "direct" &&
		school.Status == "active"
}

func schoolsByID(schools []practiqapi.SchoolInfo) map[string]practiqapi.SchoolInfo {
	index := make(map[string]practiqapi.SchoolInfo, len(schools))
	for _, school := range schools {
		index[school.ID] = school
	}
	return index
}

// List shows every Campus institution to a platform superadmin, joined with
// what practiq-be currently says about each school.
func List(repo tenantRepo.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenants, err := repo.List(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "tenant:list-error", "message": "could not list Campus institutions"})
			return
		}
		schools, err := practiq.ListAllSchools(c, c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"code": "tenant:scope-unavailable", "message": "could not reach Practiq"})
			return
		}
		index := schoolsByID(schools)

		result := make([]adminData, 0, len(tenants))
		for _, value := range tenants {
			school := index[value.SchoolID]
			result = append(result, adminData{
				ID:       value.ID,
				SchoolID: value.SchoolID,
				Name:     school.Name,
				Status:   value.Status,
				Eligible: eligible(school),
			})
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}

// Activate enables Campus for one institution.
//
// The school is verified against practiq-be rather than taken from the
// request: the caller supplies an id, and everything that decides whether it
// qualifies — that it exists, is an institution, is billed by contract and is
// still open — is read from the service that owns those facts.
func Activate(repo tenantRepo.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			SchoolID string `json:"school_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "tenant:bad-request", "message": "school_id is required"})
			return
		}

		schools, err := practiq.ListAllSchools(c, c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"code": "tenant:scope-unavailable", "message": "could not reach Practiq"})
			return
		}
		school, found := schoolsByID(schools)[input.SchoolID]
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"code": "tenant:school-not-found", "message": "no such school in Practiq"})
			return
		}
		if !eligible(school) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    "tenant:school-not-eligible",
				"message": "Campus serves active institutions billed by contract; this school is not one",
			})
			return
		}

		value, err := repo.Create(c, domain.Tenant{SchoolID: school.ID, Status: domain.TenantStatusActive})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "tenant:activate-error", "message": "could not enable Campus for this institution"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": adminData{
			ID: value.ID, SchoolID: value.SchoolID, Name: school.Name, Status: value.Status, Eligible: true,
		}})
	}
}

// SetStatus suspends, reopens or closes a Campus institution.
//
// Suspension and closure both stop access — RequireTenant admits only an
// active tenant — and neither deletes anything: the institution's courses and
// work survive so a suspension can be lifted and a closure can be looked at.
func SetStatus(repo tenantRepo.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "tenant:bad-request", "message": "status is required"})
			return
		}
		switch input.Status {
		case domain.TenantStatusActive, domain.TenantStatusSuspended, domain.TenantStatusClosed:
		default:
			c.JSON(http.StatusBadRequest, gin.H{"code": "tenant:bad-status", "message": "status must be active, suspended or closed"})
			return
		}

		existing, err := repo.Get(c, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "tenant:lookup-error", "message": "could not read the institution"})
			return
		}
		if existing == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "tenant:not-found", "message": "no such Campus institution"})
			return
		}

		// Reopening re-checks the school: one closed in Practiq, or moved onto
		// a subscription, must not come back through Campus.
		if input.Status == domain.TenantStatusActive {
			schools, err := practiq.ListAllSchools(c, c.GetHeader("Authorization"))
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"code": "tenant:scope-unavailable", "message": "could not reach Practiq"})
				return
			}
			school, found := schoolsByID(schools)[existing.SchoolID]
			if !found || !eligible(school) {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"code":    "tenant:school-not-eligible",
					"message": "the school no longer qualifies for Campus",
				})
				return
			}
		}

		value, err := repo.SetStatus(c, existing.ID, input.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "tenant:status-error", "message": "could not change the institution's status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": adminData{
			ID: value.ID, SchoolID: value.SchoolID, Status: value.Status,
		}})
	}
}
