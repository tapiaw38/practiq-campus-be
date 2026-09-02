package enrollment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucEnrollment "github.com/tapiaw38/practiq-campus-be/internal/usecases/enrollment"
)

func NewListByCourseHandler(uc ucEnrollment.ListByCourseUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		output, appErr := uc.Execute(c, userID, isSuperAdmin, courseID, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewListMineHandler(uc ucEnrollment.ListMineUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := middlewares.GetUserID(c)

		output, appErr := uc.Execute(c, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
