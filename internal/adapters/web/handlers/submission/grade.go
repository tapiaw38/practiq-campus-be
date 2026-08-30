package submission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucSubmission "github.com/tapiaw38/practiq-campus-be/internal/usecases/submission"
)

type gradeInput struct {
	Score    int    `json:"score"`
	Feedback string `json:"feedback"`
}

func NewGradeHandler(uc ucSubmission.GradeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionID := c.Param("id")
		var input gradeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		output, appErr := uc.Execute(c, userID, isSuperAdmin, submissionID, ucSubmission.GradeInput{
			Score:    input.Score,
			Feedback: input.Feedback,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
