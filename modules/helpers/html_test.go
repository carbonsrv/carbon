package helpers

import (
	"github.com/carbonsrv/carbon/ctest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	r := gin.New()
	r.GET("/string", func(c *gin.Context) {
		String(c, 200, "Hello world!")
	})

	w := ctest.Request(r, "GET", "/string")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Hello world!")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain")
}

func TestHTMLString(t *testing.T) {
	r := gin.New()
	r.GET("/htmlstring", func(c *gin.Context) {
		HTMLString(c, 200, "<h1>Hello world!</h1>")
	})

	w := ctest.Request(r, "GET", "/htmlstring")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "<h1>Hello world!</h1>")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/html")
}
