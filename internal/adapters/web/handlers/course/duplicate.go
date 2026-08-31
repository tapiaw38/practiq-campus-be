package course

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucCourse "github.com/tapiaw38/practiq-campus-be/internal/usecases/course"
)

func NewDuplicateHandler(uc ucCourse.DuplicateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, output)
	}
}
