package middleware

import (
	"fmt"
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
	fmt.Println("Flushing or nah?")
	if fw.f != nil {
		fmt.Println("Flushin!")
		fw.f.Flush()
	}
	return
}

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
	} else {
		return func(c *gin.Context) {
			handler := cgi.Handler{
				Path: path,
				Dir:  dir,
				Args: append(args, c.Request.URL.Path),
				Env:  append(append(env, "SCRIPT_FILENAME="+dir+c.Request.URL.Path), "SCRIPT_NAME="+c.Request.URL.Path),
			}
			fw := flushWriter{c.Writer, flushfields{orig_writer: c.Writer}}
			if f, ok := c.Writer.(http.Flusher); ok {
				fmt.Println("Flusher found! Yay!")
				fw.f = f
			}
			handler.ServeHTTP(fw, c.Request)
		}
	}
}
