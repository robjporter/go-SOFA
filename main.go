package main

import (
  //"time"
  //"context"
  //"strings"
	//"net/url"
  //"mime/multipart"
  "./bootstrap"
  "./middleware/identity"
  "./routes"
  //"github.com/kataras/iris"
//	"github.com/kataras/iris/core/host"
//  "github.com/kataras/iris/middleware/recover"
//  "github.com/kataras/iris/middleware/logger"
//  "github.com/prometheus/client_golang/prometheus/promhttp"
//  prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
)


func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("Awesome App", "roporter@cisco.com")
	app.Bootstrap()
	app.Configure(identity.Configure, routes.Configure)
	return app
}

func main() {
	app := newApp()
	app.Listen("secure")
}

/*
func main() {

	  app.RegisterView(iris.HTML("./views", ".html").Reload(true))

    v1 := app.Party("/api/v1")
  	{
  		v1.Post("/login", loginEndpoint)
  		v1.Post("/submit", submitEndpoint)
  		v1.Post("/read", readEndpoint)
  	}

    app.Get("/ping", func(ctx iris.Context) {
        ctx.JSON(iris.Map{
            "message": "pong",
        })
    })

    app.Get("/", indexPage)

    app.Post("/upload", iris.LimitRequestBodySize(maxSize), func(ctx iris.Context) {
        //
        // UploadFormFiles
        // uploads any number of incoming files ("multiple" property on the form input).
        //

        // The second, optional, argument
        // can be used to change a file's name based on the request,
        // at this example we will showcase how to use it
        // by prefixing the uploaded file with the current user's ip.
        ctx.UploadFormFiles("./uploads", beforeSave)
    })

	  app.Get("/metrics", iris.FromStd(promhttp.Handler()))

    // listen and serve on http://0.0.0.0:8080.
    //app.Run(iris.Addr(":8080"), iris.WithPostMaxMemory(maxSize))

    iris.RegisterOnInterrupt(func() {
        timeout := 5 * time.Second
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()
        // close all hosts
        app.Shutdown(ctx)
    })

  	target, _ := url.Parse("https://127.0.1:443")
  	go host.NewProxy(":80", target).ListenAndServe()

  	// start the server (HTTPS) on port 443, this is a blocking func
  	app.Run(iris.TLS(":443", "server.crt", "server.key"), iris.WithPostMaxMemory(maxSize), iris.WithConfiguration(iris.YAML("./config/server.yml")))
}

func indexPage(ctx iris.Context) {
  ctx.ViewData("Title", DefaultTitle)
	ctx.ViewLayout(DefaultLayout)
  if err := ctx.View("index.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func loginEndpoint(ctx iris.Context) {}
func submitEndpoint(ctx iris.Context) {}
func readEndpoint(ctx iris.Context) {}

func beforeSave(ctx iris.Context, file *multipart.FileHeader) {
    ip := ctx.RemoteAddr()
    // make sure you format the ip in a way
    // that can be used for a file name (simple case):
    ip = strings.Replace(ip, ".", "_", -1)
    ip = strings.Replace(ip, ":", "_", -1)

    // you can use the time.Now, to prefix or suffix the files
    // based on the current time as well, as an exercise.
    // i.e unixTime :=	time.Now().Unix()
    // prefix the Filename with the $IP-
    // no need for more actions, internal uploader will use this
    // name to save the file into the "./uploads" folder.
    file.Filename = ip + "-" + file.Filename
}
*/
