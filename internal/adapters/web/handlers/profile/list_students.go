package profile

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ucProfile "github.com/tapiaw38/practiq-campus-be/internal/usecases/profile"
)

func NewListStudentsHandler(uc ucProfile.ListStudentsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		output, appErr := uc.Execute(c, ucProfile.ListStudentsInput{Search: c.Query("search"), Page: page, PerPage: 20})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
