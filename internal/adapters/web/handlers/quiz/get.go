package quiz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucQuiz "github.com/tapiaw38/practiq-campus-be/internal/usecases/quiz"
)

func NewGetHandler(uc ucQuiz.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, c.Param("id"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
