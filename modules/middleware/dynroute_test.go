package middleware

import (
	"github.com/carbonsrv/carbon/ctest"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExtRoute(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/*path", ExtRoute(Plan{
		".js": Echo(200, "javascript"),
		"***": Echo(200, "not javascript"),
	}))
	Convey("Given the Extension routing route /*path", t, func() {
		Convey("When a request hits on /some.js", func() {
			w := ctest.Request(r, "GET", "/some.js")
			Convey("The Response Code should be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})
			Convey("The Body should be \"javascript\"", func() {
				So(w.Body.String(), ShouldEqual, "javascript")
			})
		})
		Convey("When a request hits on /notsome.js.php", func() {
			w := ctest.Request(r, "GET", "/notsome.js.php")
			Convey("The Response Code should be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})
			Convey("The Body should be \"not javascript\"", func() {
				So(w.Body.String(), ShouldEqual, "not javascript")
			})
		})
	})
}
