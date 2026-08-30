package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
)

type syncInput struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func NewSyncHandler(uc ucProfile.SyncUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input syncInput
		// Body is optional: a bare POST /profile with no name is still a
		// valid first sync.
		_ = c.ShouldBindJSON(&input)

		userID := middlewares.GetUserID(c)

		// profile_type comes from the token's role, never the body — the same
		// guard practiq-be uses, so a student cannot declare itself a teacher.
		profileType := "student"
		if middlewares.IsTeacher(c) {
			profileType = "teacher"
		}

		output, appErr := uc.Execute(c, ucProfile.SyncInput{
			ID:          userID,
			ProfileType: profileType,
			FullName:    input.FullName,
			Email:       input.Email,
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
