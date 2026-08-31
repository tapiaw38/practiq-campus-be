package assignment

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucAssignment "github.com/tapiaw38/practiq-campus-be/internal/usecases/assignment"
)

type createInput struct {
	SectionID       *string `json:"section_id"`
	Title           string  `json:"title" binding:"required"`
	Description     string  `json:"description"`
	DueAt           *string `json:"due_at"`
	MaxScore        int     `json:"max_score"`
	Weight          int     `json:"weight"`
	VisibleGroupID  *string `json:"visible_group_id"`
	UnlockAfterType *string `json:"unlock_after_type"`
	UnlockAfterID   *string `json:"unlock_after_id"`
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

func NewCreateHandler(uc ucAssignment.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		output, appErr := uc.Execute(c, userID, isSuperAdmin, courseID, ucAssignment.CreateInput{
			SectionID:       input.SectionID,
			Title:           input.Title,
			Description:     input.Description,
			DueAt:           parseDateTime(input.DueAt),
			MaxScore:        input.MaxScore,
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

		c.JSON(http.StatusCreated, output)
	}
}
