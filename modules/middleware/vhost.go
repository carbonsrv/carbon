package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
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
		} else if plan[hostwithoutport] != nil {
			fmt.Println("Found without port")
			plan[hostwithoutport](c)
		} else if plan["***"] != nil {
			fmt.Println("Found catchall")
			plan["***"](c)
		} else {
			fmt.Println("Found nothing")
		}
	}
}
