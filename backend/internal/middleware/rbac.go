package middleware

import (
	"net/http"

	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/JuanCJR/task-controller/internal/repository"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(permissionRepo *repository.PermissionRepository, action model.Action, module model.Module) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		mustChange, err := permissionRepo.NeedToChangePassword(c.Request.Context(), userID.(string))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking password status"})
			c.Abort()
			return
		}

		if mustChange {
			c.JSON(http.StatusForbidden, gin.H{"error": "password change required"})
			c.Abort()
			return
		}

		hasPermission, err := permissionRepo.UserHasPermission(
			c.Request.Context(),
			userID.(string),
			action,
			module,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking permissions"})
			c.Abort()
			return
		}
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
