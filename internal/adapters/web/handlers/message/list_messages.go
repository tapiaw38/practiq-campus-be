package message

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucMessage "github.com/tapiaw38/practiq-campus-be/internal/usecases/message"
)

func NewListMessagesHandler(uc ucMessage.ListMessagesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversationID := c.Param("id")
		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID, conversationID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
