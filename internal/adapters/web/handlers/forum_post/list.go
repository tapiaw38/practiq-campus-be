package forum_post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	repo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/forum_post"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucPost "github.com/tapiaw38/practiq-campus-be/internal/usecases/forum_post"
)

func NewListHandler(uc ucPost.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		options := repo.ListOptions{Limit: 20}
		if value, err := strconv.Atoi(c.Query("limit")); err == nil && value > 0 {
			options.Limit = value
		}
		if value, err := strconv.Atoi(c.Query("offset")); err == nil && value >= 0 {
			options.Offset = value
		}

		output, appErr := uc.Execute(c, userID, isSuperAdmin, threadID, options)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
