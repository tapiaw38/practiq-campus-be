package calendar_event

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucCalendar "github.com/tapiaw38/practiq-campus-be/internal/usecases/calendar_event"
)

type createInput struct {
	CourseID        *string  `json:"course_id"`
	AttendeeIDs     []string `json:"attendee_ids"`
	Title           string   `json:"title" binding:"required"`
	Description     string   `json:"description"`
	StartsAt        string   `json:"starts_at" binding:"required"`
	EndsAt          *string  `json:"ends_at"`
	AllDay          bool     `json:"all_day"`
	RecurrenceRule  string   `json:"recurrence_rule"`
	ReminderMinutes *int     `json:"reminder_minutes"`
}

func NewCreateHandler(uc ucCalendar.CreateUsecase) gin.HandlerFunc {
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
			if t, err := time.Parse(time.RFC3339, *input.EndsAt); err == nil {
				endsAt = &t
			}
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID, ucCalendar.CreateInput{
			CourseID:        input.CourseID,
			AttendeeIDs:     input.AttendeeIDs,
			Title:           input.Title,
			Description:     input.Description,
			StartsAt:        startsAt,
			EndsAt:          endsAt,
			AllDay:          input.AllDay,
			RecurrenceRule:  input.RecurrenceRule,
			ReminderMinutes: input.ReminderMinutes,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
