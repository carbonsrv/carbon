package helpers

import (
	"github.com/carbonsrv/carbon/ctest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHTMLHelpers(t *testing.T) {
	r := gin.New()
	r.GET("/string", func(c *gin.Context) {
		String(c, 200, "Hello world!")
	})

	w := ctest.Request(r, "GET", "/string")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, string(w.Body), "Hello world!")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain")
}
