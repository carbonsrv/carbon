package middleware

import (
	"github.com/carbonsrv/carbon/ctest"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIf_Regex(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r_if, _ := If_Regexp("\\.js$", Echo(200, "javascript"))
	r_ifn, _ := If_Not_Regexp("\\.js$", Echo(200, "not javascript"))
	r.GET("/*path", Combine([]func(*gin.Context){
		r_if,
		r_ifn,
	}))
	Convey("Given the match-all route /*path", t, func() {
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
