package enrollment

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/enrollment"
	"net/http"
)

func NewSelfHandler(u uc.SelfUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		o, e := u.Execute(c, middlewares.GetUserID(c), c.Param("id"))
		if e != nil {
			e.Log(c)
			c.JSON(e.StatusCode(), e)
			return
		}
		c.JSON(http.StatusCreated, o)
	}
}
