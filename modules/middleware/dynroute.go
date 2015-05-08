package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

type Plan map[string]func(*gin.Context)

// Dynamic routing based on file extension given by a map.
func ExtRoute(plan Plan) func(*gin.Context) {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		ext := filepath.Ext(path)
		fmt.Println(path)
		fmt.Println(ext)
		//found := false
		/*for ext, fn := range plan {
			if filepath.Ext(path) == ext {
				fn(c)
				//found = true
				return
			}
		}*/
		if plan[ext] != nil {
			plan[ext](c)
			return
		}
		if plan["***"] != nil {
			plan["***"](c)
			return
		}
		c.Next()
		/*if !found {
			c.Next()
		}*/
	}
}

// Old lua/static router
/*func logic_switcheroo(dir string, cfe *cache.Cache) func(*gin.Context) {
	st := staticServe.ServeCached("", staticServe.PhysFS("", true, true), cfe)
	lr := middleware.Lua(dir)
	return func(context *gin.Context) {
		file := dir + context.Request.URL.Path
		fe := cacheFileExists(file)
		if fe == true {
			if strings.HasSuffix(file, ".lua") {
				lr(context)
			} else {
				st(context)
			}
		} else {
			context.Next()
		}
	}
}*/
