package quiz_attempt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucQuizAttempt "github.com/tapiaw38/practiq-campus-be/internal/usecases/quiz_attempt"
)

func NewSaveDraftHandler(uc ucQuizAttempt.SaveDraftUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input submitInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		answers := make([]ucQuizAttempt.AnswerInput, 0, len(input.Answers))
		for _, a := range input.Answers {
			answers = append(answers, ucQuizAttempt.AnswerInput{QuestionID: a.QuestionID, AnswerText: a.AnswerText})
		}

		if appErr := uc.Execute(c, middlewares.GetUserID(c), c.Param("id"), answers); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
