package course_group

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/course_group"
)

type createInput struct {
	Name string `json:"name" binding:"required"`
}

func Create(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		group, appErr := u.Create(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), input.Name)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": group})
	}
}
