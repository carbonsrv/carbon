package glue

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGlue(t *testing.T) {
	Convey("Accessing the built glue", t, func() {
		Convey("File 'gluetest' should equal \"Hello world!\\n\"", func() {
			So(GetGlue("gluetest"), ShouldEqual, "Hello world!\n")
		})
	})
}
