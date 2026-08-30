package forum_thread

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucThread "github.com/tapiaw38/practiq-campus-be/internal/usecases/forum_thread"
)

type updateInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func NewUpdateHandler(uc ucThread.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), ucThread.UpdateInput{Title: input.Title, Description: input.Description})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
