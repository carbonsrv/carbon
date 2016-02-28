package middleware

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/cgi"
)

// To allow streaming.
type flushfields struct {
	f           http.Flusher
	orig_writer io.Writer
}
type flushWriter struct {
	http.ResponseWriter
	flushfields
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.orig_writer.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

// CGI runs a CGI app that is statically selected.
func CGI(path, dir string, args, env []string) func(*gin.Context) {
	handler := cgi.Handler{
		Path: path,
		Dir:  dir,
		Args: args,
		Env:  env,
	}
	return func(c *gin.Context) {
		fw := flushWriter{c.Writer, flushfields{orig_writer: c.Writer}}
		if f, ok := c.Writer.(http.Flusher); ok {
			fw.f = f
		}
		handler.ServeHTTP(fw, c.Request)
	}
}

// CGI_Dynamic behaves like a normal web server would do when configured to run cgi apps.
func CGI_Dynamic(path, dir string, args, env []string) func(*gin.Context) {
	if path == "" {
		return func(c *gin.Context) {
			handler := cgi.Handler{
				Path: dir + c.Request.URL.Path,
				Dir:  dir,
				Args: append(args, c.Request.URL.Path),
				Env:  append(append(env, "SCRIPT_FILENAME="+dir+c.Request.URL.Path), "SCRIPT_NAME="+c.Request.URL.Path),
			}
			fw := flushWriter{c.Writer, flushfields{orig_writer: c.Writer}}
			if f, ok := c.Writer.(http.Flusher); ok {

				fw.f = f
			}
			handler.ServeHTTP(fw, c.Request)
		}
	}
	return func(c *gin.Context) {
		handler := cgi.Handler{
			Path: path,
			Dir:  dir,
			Args: append(args, c.Request.URL.Path),
			Env:  append(append(env, "SCRIPT_FILENAME="+dir+c.Request.URL.Path), "SCRIPT_NAME="+c.Request.URL.Path),
		}
		fw := flushWriter{c.Writer, flushfields{orig_writer: c.Writer}}
		if f, ok := c.Writer.(http.Flusher); ok {
			fw.f = f
		}
		handler.ServeHTTP(fw, c.Request)
	}
}
