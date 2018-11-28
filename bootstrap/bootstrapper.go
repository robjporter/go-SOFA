package bootstrap

import (
	"time"
  "net/url"
	"github.com/kataras/iris"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/websocket"
	"github.com/kataras/iris/core/host"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
)

const (
	DefaultTitle  = "My Awesome Site"
	DefaultLayout = "layouts/layout.html"
  maxSize = 10 << 20 // 10MB
	StaticAssets = "./public/"
  SessionsDurationHour = 24
	Favicon = "favicon.ico"
  ViewDir = "./views"
  ViewsDirExtension = ".html"
  ViewsDirSharedLayout = "shared/layout.html"
  ViewsDirSharedError = "shared/error.html"
  ConfigFile = "./config/server.yml"
  SecretHash = "the-big-and-secret-fash-key-here"
  SecretHash2 = "lot-secret-of-characters-big-too"
  TLSKey = "server.key"
  TLSCert = "server.crt"
  TLSPort = "443"
  PrometheusServiceName = "serviceName"
  PrometheusURL = "/metrics"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time
	Development bool
	Sessions *sessions.Sessions
}

// New returns a new Bootstrapper.
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
		Development: true,
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

func (b *Bootstrapper) SetDevelopment(setting bool) {
	b.Development = setting
}

// SetupViews loads the templates.
func (b *Bootstrapper) SetupViews(viewsDir string) {
	if(b.Development) {
		b.RegisterView(iris.HTML(viewsDir, ViewsDirExtension).Layout(ViewsDirSharedLayout).Reload(true))
	} else {
		b.RegisterView(iris.HTML(viewsDir, ViewsDirExtension).Layout(ViewsDirSharedLayout))
	}
}

// SetupSessions initializes the sessions, optionally.
func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + b.AppName,
		Expires:  expires,
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

// SetupWebsockets prepares the websocket server.
func (b *Bootstrapper) SetupWebsockets(endpoint string, onConnection websocket.ConnectionFunc) {
	ws := websocket.New(websocket.Config{})
	ws.OnConnection(onConnection)

	b.Get(endpoint, ws.Handler())
	b.Any("/iris-ws.js", func(ctx iris.Context) {
		ctx.Write(websocket.ClientSource)
	})
}

// SetupErrorHandlers prepares the http error handlers
// `(context.StatusCodeNotSuccessful`,  which defaults to < 200 || >= 400 but you can change it).
func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		ctx.View(ViewsDirSharedError)
	})
}

// Configure accepts configurations and runs them inside the Bootstraper's context.
func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

// Bootstrap prepares our application.
//
// Returns itself.
func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews(ViewDir)
	b.SetupSessions(SessionsDurationHour*time.Hour,
		[]byte(SecretHash),
		[]byte(SecretHash2),
	)
	b.SetupErrorHandlers()

	// static files
	b.Favicon(StaticAssets + Favicon)
	b.StaticWeb(StaticAssets[1:len(StaticAssets)-1], StaticAssets)

	// middleware, after static files

  b.SetUpMiddleware()
  b.SetupDefaultRoutes()

	return b
}

func (b *Bootstrapper) SetUpMiddleware() {
	b.Use(recover.New())
	b.Use(logger.New())

	m := prometheusMiddleware.New(PrometheusServiceName)
	b.Use(m.ServeHTTP)
}

func (b *Bootstrapper) SetupDefaultRoutes() {
  b.Get(PrometheusURL, iris.FromStd(promhttp.Handler()))
}

// Listen starts the http server with the specified "addr".
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
  if(addr == "secure") {
    target, _ := url.Parse("https://127.0.1:443")
    go host.NewProxy(":8080", target).ListenAndServe()
    go host.NewProxy(":80", target).ListenAndServe()
    if(cfgs == nil) {
      b.Run(iris.TLS(":"+TLSPort, TLSCert, TLSKey), iris.WithPostMaxMemory(maxSize), iris.WithConfiguration(iris.YAML(ConfigFile)))
    } else {
      b.Run(iris.TLS(":"+TLSPort, TLSCert, TLSKey), cfgs...)
    }
  } else {
    if(cfgs == nil) {
      b.Run(iris.Addr(addr), iris.WithPostMaxMemory(maxSize), iris.WithConfiguration(iris.YAML(ConfigFile)))
    } else {
      b.Run(iris.Addr(addr), cfgs...)
    }
  }
	//
}
