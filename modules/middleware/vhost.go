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
		c_cp = c.Copy()
		host := c_cp.Request.Host
		hostwithoutport := portmatch.ReplaceAllLiteralString(host, "")
		fmt.Println(hostwithoutport)
		if plan[host] != nil {
			fmt.Println("Found with port")
			plan[host](c)
			return
		}
		if plan[hostwithoutport] != nil {
			fmt.Println("Found without port")
			plan[hostwithoutport](c)
			return
		}
		if plan["***"] != nil {
			fmt.Println("Found nothing")
			plan["***"](c)
			return
		}
	}
}
