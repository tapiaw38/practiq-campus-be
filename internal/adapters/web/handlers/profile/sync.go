package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
)

func NewSyncHandler(uc ucProfile.SyncUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := middlewares.GetUserID(c)

		profileType := "student"
		if middlewares.IsTeacher(c) {
			profileType = "teacher"
		}

		output, appErr := uc.Execute(c, ucProfile.SyncInput{
			ID:          userID,
			ProfileType: profileType,
			BearerToken: c.GetHeader("Authorization"),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
