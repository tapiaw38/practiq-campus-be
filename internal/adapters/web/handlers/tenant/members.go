package tenant

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
	ucTenant "github.com/tapiaw38/practiq-campus-be/internal/usecases/tenant"
)

func requireSchoolAdmin(c *gin.Context) bool {
	if middlewares.IsSuperAdmin(c) || tenantcontext.IsAdmin(c) {
		return true
	}
	c.JSON(http.StatusForbidden, gin.H{"code": "tenant:admin-required", "message": "school administrator access required"})
	return false
}

func SchoolMembers(uc ucTenant.ListMembersUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !requireSchoolAdmin(c) {
			return
		}
		output, appErr := uc.Execute(c, ucTenant.ListMembersInput{
			BearerToken: c.GetHeader("Authorization"),
			TenantID:    middlewares.GetTenantID(c),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func AddSchoolMember(uc ucTenant.AddMemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !requireSchoolAdmin(c) {
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

		appErr := uc.Execute(c, ucTenant.AddMemberInput{
			BearerToken: c.GetHeader("Authorization"),
			TenantID:    middlewares.GetTenantID(c),
			UserID:      strings.TrimSpace(input.UserID),
			Role:        input.Role,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func RemoveSchoolMember(uc ucTenant.RemoveMemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !requireSchoolAdmin(c) {
			return
		}

		appErr := uc.Execute(c, ucTenant.RemoveMemberInput{
			BearerToken: c.GetHeader("Authorization"),
			TenantID:    middlewares.GetTenantID(c),
			UserID:      c.Param("userID"),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
