package middleware

import (
	"github.com/gin-gonic/gin"
	"regexp"
)

// If_Written runs the handler if the content is written
func If_Written(handler func(*gin.Context)) func(*gin.Context) { // Runs handler if the content is already written.
	return func(c *gin.Context) {
		if c.Writer.Written() {
			handler(c)
		}
	}
}

// If_Regexp runs when URL matches regexp
func If_Regexp(regex string, handler func(*gin.Context)) (func(*gin.Context), error) { // Runs if the URL matches the given regexp, otherwise does nothing.
	expr, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	return func(c *gin.Context) {
		if expr.MatchString(c.Request.URL.Path) {
			handler(c)
		}
	}, nil
}

// If_Status runs if the status is the same.
func If_Status(status int, handler func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		if c.Writer.Status() == status {
			handler(c)
		}
	}
}

// And below the inverted..

// If_Not_Written is the inverted of If_Written
func If_Not_Written(handler func(*gin.Context)) func(*gin.Context) { // Runs handler if the content is already written.
	return func(c *gin.Context) {
		if !c.Writer.Written() {
			handler(c)
		}
	}
}

// If_Not_Regexp is the inverted of If_Regexp
func If_Not_Regexp(regex string, handler func(*gin.Context)) (func(*gin.Context), error) { // Runs if the URL matches the given regexp, otherwise does nothing.
	expr, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	return func(c *gin.Context) {
		if !expr.MatchString(c.Request.URL.Path) {
			handler(c)
		}
	}, nil
}

// If_Not_Status is teh inverted of If_Status
func If_Not_Status(status int, handler func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		if c.Writer.Status() != status {
			handler(c)
		}
	}
}
