package middleware

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// Authorization creates a middleware that enforces Casbin policies.
func Authorization(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user's role from the context.
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
			return
		}

		//  Define the object (URL path) and action (HTTP method).
		obj := c.Request.URL.Path // e.g., "/api/v1/asset/create"
		act := c.Request.Method   // e.g., "POST"

		// Ask Casbin to enforce the policy.
		// Check: Can 'role' perform 'act' on 'obj'?
		ok, err := enforcer.Enforce(role, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Authorization error: " + err.Error()})
			return
		}
		log.Printf("Authorization check for role=%v, obj=%s, act=%s: %v", role, obj, act, ok)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
			return
		}

		c.Next()
	}
}
