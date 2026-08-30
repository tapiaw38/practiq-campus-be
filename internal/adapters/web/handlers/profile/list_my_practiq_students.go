package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
)

func NewListMyPractiqStudentsHandler(uc ucProfile.ListMyPractiqStudentsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
