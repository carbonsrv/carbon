package middleware

import (
	"github.com/gin-gonic/gin"
	"regexp"
)

// VHOST does dynamic routing based on host given by a map.
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
		} else {
			c.Data(404, "text/plain", []byte("404 page not found"))
		}
	}
}

// VHOST_Middleware does the same as above, but for middleware and not for route handlers.
func VHOST_Middleware(plan Plan) gin.HandlerFunc {
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
		} else {
			c.Next()
		}
	}
}
