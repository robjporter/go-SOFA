package routes

import (
	"../bootstrap"
	"github.com/kataras/iris"
)

const (

  maxSize = 10 << 20 // 10MB
)

// Configure registers the necessary routes to the app.
func Configure(b *bootstrap.Bootstrapper) {
	b.Get("/", GetIndexHandler)
	b.Get("/api", GetAPIHandler)
	b.Get("/newreport", GetNewReportHandler)
	b.Get("/newmapping", GetNewMappingHandler)
	b.Post("/newmapping", PostNewMappingHandler)
	b.Get("/processreport", GetProcessReportHandler)
	b.Post("/upload", iris.LimitRequestBodySize(maxSize+1<<20), PostUpload)
}
