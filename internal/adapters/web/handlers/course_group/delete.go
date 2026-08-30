package course_group

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/course_group"
)

func Delete(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := u.Delete(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id")); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
