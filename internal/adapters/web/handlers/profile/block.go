package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
)

type blockInput struct {
	Blocked bool `json:"blocked"`
}

func NewBlockHandler(uc ucProfile.BlockUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input blockInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid request body"})
			return
		}
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.Param("id"), input.Blocked)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
