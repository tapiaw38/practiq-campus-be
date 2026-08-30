package course_group

import (
	"net/http"

	"github.com/gin-gonic/gin"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/course_group"
)

func List(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		groups, appErr := u.List(c, c.Param("id"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": groups})
	}
}
