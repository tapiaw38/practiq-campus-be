package course_material

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucMaterial "github.com/tapiaw38/practiq-campus-be/internal/usecases/course_material"
)

type createInput struct {
	AssignmentID *string `json:"assignment_id"`
	SectionID    *string `json:"section_id"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Kind         string  `json:"kind" binding:"required"`
	// URL is the value returned by POST /uploads for kind="file", or the
	// external address the teacher pasted for kind="link".
	URL string `json:"url" binding:"required"`
}

func NewCreateHandler(uc ucMaterial.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		output, appErr := uc.Execute(c, userID, isSuperAdmin, courseID, ucMaterial.CreateInput{
			AssignmentID: input.AssignmentID,
			SectionID:    input.SectionID,
			Title:        input.Title,
			Description:  input.Description,
			Kind:         input.Kind,
			URL:          input.URL,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
