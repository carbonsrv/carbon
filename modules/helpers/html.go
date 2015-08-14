package helpers

import (
	"github.com/gin-gonic/gin"
)

func String(context *gin.Context, responsecode int, text string) {
	context.Data(responsecode, "text/plain", []byte(text))
}

func HTMLString(context *gin.Context, responsecode int, html string) {
	context.Data(responsecode, "text/html", []byte(html))
}
