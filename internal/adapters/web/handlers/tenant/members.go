package tenant

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tenantRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/tenant"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func forwardMemberError(c *gin.Context, err error, code, message string) {
	var upstream *practiqapi.UpstreamError
	if errors.As(err, &upstream) && upstream.Status >= 400 && upstream.Status < 500 {
		c.JSON(upstream.Status, gin.H{"code": upstream.Code, "message": upstream.Message})
		return
	}
	c.JSON(http.StatusBadGateway, gin.H{"code": code, "message": message})
}

// SchoolMembers exposes membership only after RequireTenant has bound the
// request. The selected Campus tenant determines school_id; client input never
// chooses a school, so an institution admin cannot pivot into another one.
func SchoolMembers(repo tenantRepo.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !middlewares.IsSuperAdmin(c) && !tenantcontext.IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"code": "tenant:admin-required", "message": "school administrator access required"})
			return
		}
		value, err := repo.Get(c, middlewares.GetTenantID(c))
		if err != nil || value == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "tenant:not-found", "message": "institution not found"})
			return
		}
		members, err := practiq.ListSchoolMembers(c, c.GetHeader("Authorization"), value.SchoolID)
		if err != nil {
			forwardMemberError(c, err, "tenant:members-unavailable", "could not load institution members")
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": members})
	}
}

func AddSchoolMember(repo tenantRepo.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !middlewares.IsSuperAdmin(c) && !tenantcontext.IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"code": "tenant:admin-required", "message": "school administrator access required"})
			return
		}
		var input struct {
			UserID string `json:"user_id"`
			Role   string `json:"role"`
		}
		if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.UserID) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "tenant:bad-member", "message": "user_id is required"})
			return
		}
		if input.Role != "admin" && input.Role != "teacher" && input.Role != "student" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "tenant:bad-role", "message": "role must be admin, teacher, or student"})
			return
		}
		value, err := repo.Get(c, middlewares.GetTenantID(c))
		if err != nil || value == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "tenant:not-found", "message": "institution not found"})
			return
		}
		if err := practiq.AddSchoolMember(c, c.GetHeader("Authorization"), value.SchoolID, strings.TrimSpace(input.UserID), input.Role); err != nil {
			forwardMemberError(c, err, "tenant:add-member-error", "could not add member to institution")
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func RemoveSchoolMember(repo tenantRepo.Repository, practiq practiqapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !middlewares.IsSuperAdmin(c) && !tenantcontext.IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"code": "tenant:admin-required", "message": "school administrator access required"})
			return
		}
		value, err := repo.Get(c, middlewares.GetTenantID(c))
		if err != nil || value == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "tenant:not-found", "message": "institution not found"})
			return
		}
		if err := practiq.RemoveSchoolMember(c, c.GetHeader("Authorization"), value.SchoolID, c.Param("userID")); err != nil {
			forwardMemberError(c, err, "tenant:remove-member-error", "could not remove member from institution")
			return
		}
		c.Status(http.StatusNoContent)
	}
}
