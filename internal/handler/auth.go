package handler

import (
	"net/http"
	"strings"

	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"github.com/ZhangPengPaul/LightOTA/internal/repository"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Authorization header missing"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid authorization format"})
			return
		}

		apiKey := parts[1]
		tenant, err := repo.FindByAPIKey(apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid API key"})
			return
		}

		c.Set("tenant", tenant)
		c.Next()
	}
}

func GetTenant(c *gin.Context) *model.Tenant {
	tenant, exists := c.Get("tenant")
	if !exists {
		return nil
	}
	return tenant.(*model.Tenant)
}

func DeviceAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Authorization header missing"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid authorization format"})
			return
		}

		token := parts[1]
		c.Set("device_token", token)
		c.Next()
	}
}
