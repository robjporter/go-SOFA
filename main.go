package main

import (
  "./bootstrap"
  "./middleware/identity"
  "./routes"
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
