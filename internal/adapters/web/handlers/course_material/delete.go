package course_material

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucMaterial "github.com/tapiaw38/practiq-campus-be/internal/usecases/course_material"
)

func NewDeleteHandler(uc ucMaterial.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		materialID := c.Param("id")
		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		if appErr := uc.Execute(c, userID, isSuperAdmin, materialID); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
