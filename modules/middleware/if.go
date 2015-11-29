package middleware

import (
	"github.com/gin-gonic/gin"
	"regexp"
)

func If_Written(handler func(*gin.Context)) func(*gin.Context) { // Runs handler if the content is already written.
	return func(c *gin.Context) {
		if c.Writer.Written() {
			handler(c)
		}
	}
}

func If_Regexp(regex string, handler func(*gin.Context)) (func(*gin.Context), error) { // Runs if the URL matches the given regexp, otherwise does nothing.
	expr, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	return func(c *gin.Context) {
		if regexp.MatchString(expr) {
			handler(c)
		}
	}, nil
}

func If_Status(status int, handler func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		if c.Writer.Status() == status {
			handler(c)
		}
	}
}

// And below the inverted..

func If_Not_Written(handler func(*gin.Context)) func(*gin.Context) { // Runs handler if the content is already written.
	return func(c *gin.Context) {
		if !c.Writer.Written() {
			handler(c)
		}
	}
}

func If_Not_Regexp(regex string, handler func(*gin.Context)) (func(*gin.Context), error) { // Runs if the URL matches the given regexp, otherwise does nothing.
	expr, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	return func(c *gin.Context) {
		if !regexp.MatchString(expr) {
			handler(c)
		}
	}, nil
}

func If_Not_Status(status int, handler func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		if c.Writer.Status() != status {
			handler(c)
		}
	}
}
