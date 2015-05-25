package helpers

import (
	"github.com/gin-gonic/gin"
)

func HTMLString(context *gin.Context, responsecode int, html string) {
	context.Data(responsecode, "text/html", []byte(html))
}
