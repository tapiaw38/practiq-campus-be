package quiz_attempt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucQuizAttempt "github.com/tapiaw38/practiq-campus-be/internal/usecases/quiz_attempt"
)

func NewStartHandler(uc ucQuizAttempt.StartUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.Param("id"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, output)
	}
}
