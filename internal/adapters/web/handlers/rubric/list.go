package rubric

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/rubric"
)

func List(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		x, e := u.List(c, c.Param("id"), middlewares.GetUserID(c), middlewares.IsSuperAdmin(c))
		if e != nil {
			c.JSON(e.StatusCode(), e)
			return
		}
		c.JSON(200, gin.H{"data": x})
	}
}
