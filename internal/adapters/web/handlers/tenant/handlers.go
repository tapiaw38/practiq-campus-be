package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/middlewares"
	ucTenant "github.com/tapiaw38/practiq-campus-be/internal/usecases/tenant"
)

func Mine(uc ucTenant.ListMineUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, ucTenant.ListMineInput{
			BearerToken:  c.GetHeader("Authorization"),
			IsSuperAdmin: middlewares.IsSuperAdmin(c),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
