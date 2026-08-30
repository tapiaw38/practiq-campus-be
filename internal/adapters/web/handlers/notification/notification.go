package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/notification"
	"net/http"
)

func List(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		x, e := u.List(c, middlewares.GetUserID(c))
		if e != nil {
			c.JSON(e.StatusCode(), e)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": x})
	}
}
func Read(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if e := u.Read(c, c.Param("id"), middlewares.GetUserID(c)); e != nil {
			c.JSON(e.StatusCode(), e)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
