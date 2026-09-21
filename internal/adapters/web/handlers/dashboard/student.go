package dashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/dashboard"
	"net/http"
)

func Student(u uc.StudentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		out, err := u.Execute(c, middlewares.GetUserID(c))
		if err != nil {
			c.JSON(err.StatusCode(), err)
			return
		}
		c.JSON(http.StatusOK, out)
	}
}
