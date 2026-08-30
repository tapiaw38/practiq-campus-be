package submission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucSubmission "github.com/tapiaw38/practiq-campus-be/internal/usecases/submission"
)

func NewGetMineHandler(uc ucSubmission.GetMineUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentID := c.Param("id")
		userID := middlewares.GetUserID(c)

		output, appErr := uc.Execute(c, userID, assignmentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
