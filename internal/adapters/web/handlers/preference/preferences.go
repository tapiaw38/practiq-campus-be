package preference

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucPreference "github.com/tapiaw38/practiq-campus-be/internal/usecases/preference"
)

func NewGetHandler(uc ucPreference.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.Param("scope"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

type updateInput struct {
	Settings json.RawMessage `json:"settings"`
}

func NewUpdateHandler(uc ucPreference.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
			return
		}
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.Param("scope"), ucPreference.UpdateInput{Settings: input.Settings})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
