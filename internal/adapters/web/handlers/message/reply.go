package message

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucMessage "github.com/tapiaw38/practiq-campus-be/internal/usecases/message"
)

type replyInput struct {
	Body string `json:"body" binding:"required"`
}

func NewReplyHandler(uc ucMessage.ReplyUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conversationID := c.Param("id")
		var input replyInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID, conversationID, ucMessage.ReplyInput{Body: input.Body})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
