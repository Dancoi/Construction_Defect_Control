package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that allows access only to users whose role
// is in the provided list. It expects the JWTAuthMiddleware to set "role" in context.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		rv, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "error", "error": "forbidden"})
			return
		}
		role, _ := rv.(string)
		if _, ok := allowed[role]; ok {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "error", "error": "forbidden"})
	}
}

// RequireOwnerOrRole checks ownership through a simple check function or allows the provided roles.
// ownerCheck should return true if current user is owner for requested resource.
func RequireOwnerOrRole(ownerCheck func(c *gin.Context) bool, roles ...string) gin.HandlerFunc {
	roleMw := RequireRole(roles...)
	return func(c *gin.Context) {
		if ownerCheck != nil && ownerCheck(c) {
			c.Next()
			return
		}
		// delegate to role middleware
		roleMw(c)
	}
}
