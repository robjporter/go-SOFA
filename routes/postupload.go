package routes

import (
  "io"
  "os"
	"github.com/kataras/iris"
)

const (
  uploadsDir = "./uploads/reports"
)

// GetIndexHandler handles the GET: /
func PostUpload(ctx iris.Context) {

  file, info, err := ctx.FormFile("file")
  if err != nil {
      ctx.StatusCode(iris.StatusInternalServerError)
      ctx.Application().Logger().Warnf("Error while uploading: %v", err.Error())
      return
  }

  defer file.Close()
  fname := info.Filename

  // Create a file with the same name
  // assuming that you have a folder named 'uploads'
  out, err := os.OpenFile(uploadsDir+"/"+fname,
      os.O_WRONLY|os.O_CREATE, 0666)

  if err != nil {
      ctx.StatusCode(iris.StatusInternalServerError)
      ctx.Application().Logger().Warnf("Error while preparing the new file: %v", err.Error())
      return
  }
  defer out.Close()

  io.Copy(out, file)

  ctx.Redirect("/", iris.StatusSeeOther)
}
