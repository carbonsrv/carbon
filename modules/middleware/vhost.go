package middleware

import (
	"github.com/gin-gonic/gin"
	"regexp"
)

// Dynamic routing based on host given by a map.
func VHOST(plan Plan) func(*gin.Context) {
	portmatch := regexp.MustCompile(":.*$")
	return func(c *gin.Context) {
		host := c.Request.Host
		if plan[host] != nil {
			plan[host](c)
			return
		}
		hostwithoutport := portmatch.ReplaceAllLiteralString(host, "")
		if plan[hostwithoutport] != nil {
			plan[hostwithoutport](c)
			return
		}
		if plan["***"] != nil {
			plan["***"](c)
			return
		}
		c.Next()
	}
}
