package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
)

func Authentication(jwtService contract.IJWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}
		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		ctx.Set("userID", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Set("department", claims.Department)
		ctx.Set("clearanceLevel", claims.ClearanceLevel)
		//
		log.Println("user context:", ctx.GetString("userID"), ctx.GetString("role"), ctx.GetString("department"), "wow",ctx.GetInt("clearanceLevel"), "hello", claims.ClearanceLevel)
		ctx.Next()
	}
}
