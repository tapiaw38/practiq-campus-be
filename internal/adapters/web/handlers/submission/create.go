package submission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucSubmission "github.com/tapiaw38/practiq-campus-be/internal/usecases/submission"
)

type createInput struct {
	Content string `json:"content" binding:"required"`
}

func NewCreateHandler(uc ucSubmission.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentID := c.Param("id")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID, assignmentID, ucSubmission.CreateInput{Content: input.Content})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
