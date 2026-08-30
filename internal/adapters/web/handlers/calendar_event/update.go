package calendar_event

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucCalendar "github.com/tapiaw38/practiq-campus-be/internal/usecases/calendar_event"
)

func NewUpdateHandler(uc ucCalendar.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		startsAt, err := time.Parse(time.RFC3339, input.StartsAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid starts_at"})
			return
		}
		var endsAt *time.Time
		if input.EndsAt != nil && *input.EndsAt != "" {
			if parsed, err := time.Parse(time.RFC3339, *input.EndsAt); err == nil {
				endsAt = &parsed
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid ends_at"})
				return
			}
		}
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.Param("id"), ucCalendar.CreateInput{CourseID: input.CourseID, AttendeeIDs: input.AttendeeIDs, Title: input.Title, Description: input.Description, StartsAt: startsAt, EndsAt: endsAt, AllDay: input.AllDay, RecurrenceRule: input.RecurrenceRule, ReminderMinutes: input.ReminderMinutes})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
