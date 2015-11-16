package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/carbonsrv/carbon/modules/helpers"
)

func Echo(status int, s string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.String(status, s)
	}
}
func EchoHTML(status int, s string) func(*gin.Context) {
	return func(c *gin.Context) {
		helpers.HTMLString(c, status, s)
	}
}
