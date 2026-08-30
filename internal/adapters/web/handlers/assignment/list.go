package assignment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucAssignment "github.com/tapiaw38/practiq-campus-be/internal/usecases/assignment"
)

func NewListHandler(uc ucAssignment.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		output, appErr := uc.Execute(c, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
