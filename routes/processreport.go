
package routes

import (
	"github.com/kataras/iris"
)

// GetIndexHandler handles the GET: /
func GetProcessReportHandler(ctx iris.Context) {
	ctx.ViewData("Title", "Index Page")
	ctx.View("processreport.html")
}
