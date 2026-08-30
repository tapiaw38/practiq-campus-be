package assignment

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	uc "github.com/tapiaw38/practiq-campus-be/internal/usecases/assignment"
	"net/http"
)

func NewUpdateHandler(u uc.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in createInput
		if c.ShouldBindJSON(&in) != nil {
			c.JSON(400, gin.H{"message": "invalid body"})
			return
		}
		o, e := u.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), uc.UpdateInput{SectionID: in.SectionID, Title: in.Title, Description: in.Description, DueAt: parseDateTime(in.DueAt), MaxScore: in.MaxScore})
		if e != nil {
			e.Log(c)
			c.JSON(e.StatusCode(), e)
			return
		}
		c.JSON(200, o)
	}
}
func NewDeleteHandler(u uc.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if e := u.Execute(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id")); e != nil {
			e.Log(c)
			c.JSON(e.StatusCode(), e)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
