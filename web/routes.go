package main

import (
	"github.com/hoisie/web"
)

func RegisterRoutes() {

  // DEFAULT ROUTE
  web.Get("/", func(ctx * web.Context) {
		ctx.Write(toJson( "let thy object decend as if it were calescent" ));
  })

}
