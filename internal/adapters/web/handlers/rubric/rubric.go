package rubric

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/rubric"
	"net/http"
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
func Replace(u uc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var x struct {
			Criteria []uc.Criterion `json:"criteria"`
		}
		if c.ShouldBindJSON(&x) != nil {
			c.JSON(400, gin.H{"message": "invalid rubric"})
			return
		}
		if e := u.Replace(c, middlewares.GetUserID(c), c.Param("id"), middlewares.IsSuperAdmin(c), x.Criteria); e != nil {
			c.JSON(e.StatusCode(), e)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
