
package routes

import (
	"github.com/kataras/iris"
)

// GetIndexHandler handles the GET: /
func GetNewReportHandler(ctx iris.Context) {
	ctx.ViewData("Title", "Index Page")
	ctx.View("newreport.html")
}
