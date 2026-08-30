package enrollment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucEnrollment "github.com/tapiaw38/practiq-campus-be/internal/usecases/enrollment"
)

func NewDeleteHandler(uc ucEnrollment.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		enrollmentID := c.Param("id")
		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		if appErr := uc.Execute(c, userID, isSuperAdmin, enrollmentID); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
