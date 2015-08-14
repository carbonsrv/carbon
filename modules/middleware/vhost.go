package middleware

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"fmt"
)

// Dynamic routing based on host given by a map.
func VHOST(plan Plan) func(*gin.Context) {
	portmatch := regexp.MustCompile(":.*$")
	return func(c *gin.Context) {
		host := c.Request.Host
		hostwithoutport := portmatch.ReplaceAllLiteralString(host, "")
		fmt.Println(hostwithoutport)

		if plan[host] != nil {
			fmt.Println("Found")
			plan[host](c)
			return
		}
		if plan[hostwithoutport] != nil {
			fmt.Println("Found without port")
			plan[hostwithoutport](c)
			return
		}
		if plan["***"] != nil {
			fmt.Println("Found catchall")
			plan["***"](c)
			return
		}
		fmt.Println("Found nothing")
		c.Next()
	}
}
