package helpers

import (
	"github.com/gin-gonic/gin"
)

// String returns string as text/plain
func String(context *gin.Context, responsecode int, text string) {
	context.Data(responsecode, "text/plain", []byte(text))
}

// HTMLString does the same as above, but with text/html
func HTMLString(context *gin.Context, responsecode int, html string) {
	context.Data(responsecode, "text/html", []byte(html))
}
