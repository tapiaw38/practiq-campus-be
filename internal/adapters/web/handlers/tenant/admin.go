package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	ucTenant "github.com/tapiaw38/practiq-campus-be/internal/usecases/tenant"
)

func List(uc ucTenant.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, ucTenant.ListInput{BearerToken: c.GetHeader("Authorization")})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func Activate(uc ucTenant.ActivateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			SchoolID string `json:"school_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "tenant:bad-request", "message": "school_id is required"})
			return
		}

		output, appErr := uc.Execute(c, ucTenant.ActivateInput{BearerToken: c.GetHeader("Authorization"), SchoolID: input.SchoolID})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, output)
	}
}

func SetStatus(uc ucTenant.SetStatusUsecase) gin.HandlerFunc {
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

		output, appErr := uc.Execute(c, ucTenant.SetStatusInput{
			BearerToken: c.GetHeader("Authorization"),
			TenantID:    c.Param("id"),
			Status:      input.Status,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func EligibleSchools(uc ucTenant.ListEligibleSchoolsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, ucTenant.ListEligibleSchoolsInput{BearerToken: c.GetHeader("Authorization")})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
