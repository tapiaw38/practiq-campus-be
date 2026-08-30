package enrollment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucEnrollment "github.com/tapiaw38/practiq-campus-be/internal/usecases/enrollment"
)

type createInput struct {
	UserID         string `json:"user_id" binding:"required"`
	EnrollmentRole string `json:"enrollment_role"`
}

func NewCreateHandler(uc ucEnrollment.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		output, appErr := uc.Execute(c, userID, isSuperAdmin, courseID, ucEnrollment.CreateInput{
			UserID:         input.UserID,
			EnrollmentRole: input.EnrollmentRole,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
