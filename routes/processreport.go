
package routes

import (
	"fmt"
	"io/ioutil"
	"github.com/kataras/iris"
)

// GetIndexHandler handles the GET: /
func GetProcessReportHandler(ctx iris.Context) {
	var reports = []string{}
	files, err := ioutil.ReadDir("./uploads/reports")
  if err != nil {
      fmt.Println(err)
  }

  for _, f := range files {
		if(!f.IsDir()) {
				reports = append(reports,removeExtension(f.Name()))
		}
  }
	ctx.ViewData("Title", "Index Page")
	ctx.ViewData("Reports",reports)
  if err := ctx.View("processreport.html"); err != nil {
    ctx.Application().Logger().Infof(err.Error())
  }
}
