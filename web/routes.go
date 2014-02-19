package main

import (
  "net/http"
	"github.com/hoisie/web"
)

func RegisterRoutes() {

  // DEFAULT ROUTE
  web.Get("/", func(ctx * web.Context) {
    http.ServeFile(ctx, ctx.Request, "./static/landing.html")
  })
 

}
