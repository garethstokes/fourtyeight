package main

import (
  "net/http"
	"github.com/hoisie/web"
	"github.com/garethstokes/fourtyeight/personal"
)

func RegisterRoutes() {

  // DEFAULT ROUTE
  web.Get("/", func(ctx * web.Context) {
    http.ServeFile(ctx, ctx.Request, "./static/rain.html")
  })

  // Take a email down for the waiting list
  web.Post("/waitinglist", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true);

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    err := p.SaveToWaitingList( ctx.Params["email"] )
    if err != nil {
      apiError( ctx, "a problem happened saving to disk." )
      return
    }

    ok( ctx, "" )
  })

}
