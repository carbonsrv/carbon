package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http/cgi"
)

func CGI(path, dir string, args, env []string) func(*gin.Context) {
	handler := cgi.Handler{
		Path: path,
		Dir:  dir,
		Args: args,
		Env:  env,
	}
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
