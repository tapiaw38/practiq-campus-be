package course

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucCourse "github.com/tapiaw38/practiq-campus-be/internal/usecases/course"
)

func NewListHandler(uc ucCourse.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := middlewares.GetUserID(c)
		isTeacher := middlewares.IsTeacher(c)

		output, appErr := uc.Execute(c, userID, isTeacher)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
