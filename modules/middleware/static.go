package middleware

import (
	"github.com/carbonsrv/carbon/modules/helpers"
	"github.com/gin-gonic/gin"
)

// Echo the string back as text
func Echo(status int, s string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.String(status, s)
	}
}

// EchoHTML does the same as above, but for html
func EchoHTML(status int, s string) func(*gin.Context) {
	return func(c *gin.Context) {
		helpers.HTMLString(c, status, s)
	}
}
