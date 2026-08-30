package course

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucCourse "github.com/tapiaw38/practiq-campus-be/internal/usecases/course"
)

type updateInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func NewUpdateHandler(uc ucCourse.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		output, appErr := uc.Execute(c, userID, isSuperAdmin, courseID, ucCourse.UpdateInput{
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
			StartDate:   parseDate(input.StartDate),
			EndDate:     parseDate(input.EndDate),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
