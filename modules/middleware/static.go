package middleware

import (
	"github.com/gin-gonic/gin"
)

func Echo(status int, s string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.String(status, s)
	}
}
func EchoHTML(status int, s string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.HTMLString(status, s)
	}
}
