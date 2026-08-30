package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
)

func NewGetMeHandler(uc ucProfile.GetUsecase) gin.HandlerFunc {
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
