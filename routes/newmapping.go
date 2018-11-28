
package routes

import (
	"fmt"
  "io/ioutil"
	"github.com/kataras/iris"
)

type ReportClass struct {
    Name string
    Value string
}

// GetIndexHandler handles the GET: /
func GetNewMappingHandler(ctx iris.Context) {
	var reports = []string{}
	files, err := ioutil.ReadDir("./uploads/reports")
  if err != nil {
      fmt.Println(err)
  }

  for _, f := range files {
		if(!f.IsDir()) {
				reports = append(reports,f.Name())
		}
  }
	ctx.ViewData("Title", "Index Page")
	ctx.ViewData("Reports",reports)
	if err := ctx.View("newmapping.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}
