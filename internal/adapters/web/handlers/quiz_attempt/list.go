package quiz_attempt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucQuizAttempt "github.com/tapiaw38/practiq-campus-be/internal/usecases/quiz_attempt"
)

func NewListMineHandler(uc ucQuizAttempt.ListMineUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.Param("id"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewListByQuizHandler(uc ucQuizAttempt.ListByQuizUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
