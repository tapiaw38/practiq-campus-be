package quiz

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucQuiz "github.com/tapiaw38/practiq-campus-be/internal/usecases/quiz"
)

type createInput struct {
	SectionID      *string `json:"section_id"`
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"`
	TimeLimitSecs  *int    `json:"time_limit_secs"`
	MaxAttempts    int     `json:"max_attempts"`
	ScheduledAt    *string `json:"scheduled_at"`
	AvailableUntil *string `json:"available_until"`
}

func parseDateTime(v *string) *time.Time {
	if v == nil || *v == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil
	}
	return &t
}

func NewCreateHandler(uc ucQuiz.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		output, appErr := uc.Execute(c, userID, isSuperAdmin, courseID, ucQuiz.CreateInput{
			SectionID:      input.SectionID,
			Title:          input.Title,
			Description:    input.Description,
			TimeLimitSecs:  input.TimeLimitSecs,
			MaxAttempts:    input.MaxAttempts,
			ScheduledAt:    parseDateTime(input.ScheduledAt),
			AvailableUntil: parseDateTime(input.AvailableUntil),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
