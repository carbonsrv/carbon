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
		hostwithoutport := portmatch.ReplaceAllLiteralString(host, "")
		if plan[host] != nil {
			plan[host](c)
		} else if plan[hostwithoutport] != nil {
			plan[hostwithoutport](c)
		} else if plan["***"] != nil {
			plan["***"](c)
		}
	}
}
