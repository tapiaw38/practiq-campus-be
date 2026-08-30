package quiz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucQuiz "github.com/tapiaw38/practiq-campus-be/internal/usecases/quiz"
)

func NewListQuestionsHandler(uc ucQuiz.QuestionsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, appErr := uc.List(c, c.Param("id"), middlewares.GetUserID(c), middlewares.IsSuperAdmin(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}

type questionInput struct {
	Type          string   `json:"type" binding:"required"`
	Statement     string   `json:"statement" binding:"required"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer" binding:"required"`
	Points        int      `json:"points"`
}

type replaceQuestionsInput struct {
	Questions []questionInput `json:"questions"`
}

func NewReplaceQuestionsHandler(uc ucQuiz.QuestionsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input replaceQuestionsInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		questions := make([]ucQuiz.QuestionInput, 0, len(input.Questions))
		for _, q := range input.Questions {
			questions = append(questions, ucQuiz.QuestionInput{
				Type:          q.Type,
				Statement:     q.Statement,
				Options:       q.Options,
				CorrectAnswer: q.CorrectAnswer,
				Points:        q.Points,
			})
		}

		if appErr := uc.Replace(c, c.Param("id"), middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), questions); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
