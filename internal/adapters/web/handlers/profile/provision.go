package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
)

type provisionInput struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func NewProvisionStudentHandler(uc ucProfile.CreateOrSyncUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input provisionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucProfile.CreateOrSyncUserInput{
			BearerToken: c.GetHeader("Authorization"),
			Email:       input.Email,
			FirstName:   input.FirstName,
			LastName:    input.LastName,
			Password:    input.Password,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
