package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	profileRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/profile"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/auth"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/revocation"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func AuthMiddleware(profiles profileRepo.Repository, checker *revocation.Checker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "common:unauthorized", "message": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "common:unauthorized", "message": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "common:unauthorized", "message": "invalid token"})
			c.Abort()
			return
		}
		if checker != nil && !checker.Current(c.Request.Context(), authHeader, claims.UserID, claims.TokenVersion) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "common:unauthorized", "message": "session ended, please sign in again"})
			c.Abort()
			return
		}
		profile, profileErr := profiles.Get(c, claims.UserID)
		if profileErr != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"code": "profile:unavailable", "message": "could not verify account status"})
			c.Abort()
			return
		}
		if profile != nil && profile.IsBlocked {
			c.JSON(http.StatusForbidden, gin.H{"code": "profile:blocked", "message": "account is blocked in Campus"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	userID, _ := c.Get("userID")
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}

func GetRoles(c *gin.Context) []auth.RoleClaim {
	roles, _ := c.Get("roles")
	if value, ok := roles.([]auth.RoleClaim); ok {
		return value
	}
	return nil
}

func HasRole(c *gin.Context, expected ...string) bool {
	roles := GetRoles(c)
	for _, role := range roles {
		for _, exp := range expected {
			if role.Name == exp {
				return true
			}
		}
	}
	return false
}

const (
	RoleStudent    = "user"
	RoleTeacher    = "admin"
	RoleSuperAdmin = "superadmin"
)

func IsSuperAdmin(c *gin.Context) bool {
	return HasRole(c, RoleSuperAdmin)
}

func IsTeacher(c *gin.Context) bool {
	return HasRole(c, RoleTeacher, RoleSuperAdmin)
}

func RequireRoles(expected ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasRole(c, expected...) {
			c.JSON(http.StatusForbidden, gin.H{"code": "common:forbidden", "message": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireTenantStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsSuperAdmin(c) || tenantcontext.IsAdmin(c) || tenantcontext.Role(c) == tenantcontext.RoleTeacher {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"code": "tenant:staff-required", "message": "teacher or institution administrator access required"})
		c.Abort()
	}
}
