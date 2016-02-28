package middleware

import (
	"github.com/gin-gonic/gin"
)

// Combine middlewares
func Combine(middlewares []func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		for _, middleware := range middlewares {
			middleware(c)
		}
	}
}
