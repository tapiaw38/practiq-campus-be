package quiz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucQuiz "github.com/tapiaw38/practiq-campus-be/internal/usecases/quiz"
)

func NewListHandler(uc ucQuiz.ListUsecase) gin.HandlerFunc {
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
