package course

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucCourse "github.com/tapiaw38/practiq-campus-be/internal/usecases/course"
)

type createInput struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func parseDate(v *string) *time.Time {
	if v == nil || *v == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *v)
	if err != nil {
		return nil
	}
	return &t
}

func NewCreateHandler(uc ucCourse.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID, ucCourse.CreateInput{
			Title:       input.Title,
			Description: input.Description,
			StartDate:   parseDate(input.StartDate),
			EndDate:     parseDate(input.EndDate),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
