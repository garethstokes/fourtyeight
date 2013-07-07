package main

import (
	"github.com/garethstokes/web"
	"github.com/garethstokes/fourtyeight/personal"
)

func RegisterRoutes() {

  // DEFAULT ROUTE
  web.Get("/", func(ctx * web.Context) {
		ctx.Write(toJson( "let thy object decend as if it were calescent" ));
  })

  web.Get("/users/([A-Za-z0-9]+)", func(ctx * web.Context, val string) {
    personal := personal.Store()
    defer personal.CloseSession()

    person, error := personal.FindByName( val )
    if error != nil {
      apiError( ctx, error.Error() )
      return
    }

    return ctx.Write( tojson( person ) )
  })
}
