package course_group

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/course_group"
)

type memberInput struct {
	UserID string `json:"user_id" binding:"required"`
}

func AddMember(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input memberInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		if appErr := u.AddMember(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), input.UserID); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func RemoveMember(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := u.RemoveMember(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), c.Param("userId")); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
