package middleware

import (
	"fmt"
	"net/http"

	"example.com/defect-control-system/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// JWTAuthMiddleware validates Authorization Bearer token and sets user ID in context
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "missing authorization"})
			return
		}
		// expect "Bearer <token>"
		var tokenStr string
		_, err := fmt.Sscanf(auth, "Bearer %s", &tokenStr)
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "invalid authorization header"})
			return
		}
		secret := viper.GetString("jwt.secret")
		claims, err := utils.ParseJWT(secret, tokenStr)
		if err != nil {
			// log header and underlying parse error for debugging
			fmt.Printf("JWT middleware: incoming Authorization header='%s', parse error=%v\n", auth, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "error": fmt.Sprintf("invalid token: %v", err)})
			return
		}
		// put user id and role into context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
