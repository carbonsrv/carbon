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

func CGI_Dynamic(path, dir string, args, env []string) func(*gin.Context) {
	if path == "" {
		return func(c *gin.Context) {
			handler := cgi.Handler{
				Path: dir + c.Request.URL.Path,
				Dir:  dir,
				Args: append(args, c.Request.URL.Path),
				Env:  append(append(env, "SCRIPT_FILENAME="+dir+c.Request.URL.Path), "SCRIPT_NAME="+c.Request.URL.Path),
			}
			handler.ServeHTTP(c.Writer, c.Request)
		}
	} else {
		return func(c *gin.Context) {
			handler := cgi.Handler{
				Path: path,
				Dir:  dir,
				Args: append(args, c.Request.URL.Path),
				Env:  append(append(env, "SCRIPT_FILENAME="+dir+c.Request.URL.Path), "SCRIPT_NAME="+c.Request.URL.Path),
			}
			handler.ServeHTTP(c.Writer, c.Request)
		}
	}
}
