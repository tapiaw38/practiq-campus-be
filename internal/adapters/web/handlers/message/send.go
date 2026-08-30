package message

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucMessage "github.com/tapiaw38/practiq-campus-be/internal/usecases/message"
)

type sendInput struct {
	ToUserID string `json:"to_user_id"`
	ToEmail  string `json:"to_email"`
	Body     string `json:"body" binding:"required"`
}

func NewSendHandler(uc ucMessage.SendUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input sendInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		if input.ToUserID == "" && input.ToEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "to_user_id or to_email is required"})
			return
		}

		userID := middlewares.GetUserID(c)
		output, appErr := uc.Execute(c, userID, ucMessage.SendInput{
			ToUserID: input.ToUserID,
			ToEmail:  input.ToEmail,
			Body:     input.Body,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
