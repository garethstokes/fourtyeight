package main

import (
	"github.com/garethstokes/web"
	//"github.com/garethstokes/fourtyeight/personal"
)

func RegisterRoutes() {

  // DEFAULT ROUTE
  web.Get("/", func(ctx * web.Context) {
		ctx.Write(toJson( "let thy object decend as if it were calescent" ));
  })

  //web.Get("/users/([A-Za-z0-9]+)", func(ctx * web.Context, val string) {
  //})
}
