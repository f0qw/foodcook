package middleware

import (
	"net/http"

	"foodcook/internal/domain/models"
	"foodcook/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

// RootMiddleware 检查用户是否为 root 用户
func RootMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			c.Abort()
			return
		}

		// 从数据库获取用户信息
		var user models.User
		if err := database.GetDB().First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 检查用户角色
		if user.Role != models.RoleRoot {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要 root 权限"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", user)
		c.Next()
	}
}
