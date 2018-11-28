package routes

import (
	"github.com/kataras/iris"
)

// GetIndexHandler handles the GET: /
func GetAPIHandler(ctx iris.Context) {
	ctx.ViewData("Title", "Index Page")
}
