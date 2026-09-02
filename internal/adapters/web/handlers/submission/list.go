package submission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucSubmission "github.com/tapiaw38/practiq-campus-be/internal/usecases/submission"
)

func NewListByAssignmentHandler(uc ucSubmission.ListByAssignmentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentID := c.Param("id")
		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		output, appErr := uc.Execute(c, userID, isSuperAdmin, assignmentID, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
