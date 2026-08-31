package quiz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucQuiz "github.com/tapiaw38/practiq-campus-be/internal/usecases/quiz"
)

func NewUpdateHandler(uc ucQuiz.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), ucQuiz.UpdateInput{
			SectionID:       input.SectionID,
			Title:           input.Title,
			Description:     input.Description,
			TimeLimitSecs:   input.TimeLimitSecs,
			MaxAttempts:     input.MaxAttempts,
			ScheduledAt:     parseDateTime(input.ScheduledAt),
			AvailableUntil:  parseDateTime(input.AvailableUntil),
			Weight:          input.Weight,
			VisibleGroupID:  input.VisibleGroupID,
			UnlockAfterType: input.UnlockAfterType,
			UnlockAfterID:   input.UnlockAfterID,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewDeleteHandler(uc ucQuiz.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id")); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
