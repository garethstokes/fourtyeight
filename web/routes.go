package main

import (
	"github.com/garethstokes/web"
	"github.com/garethstokes/fourtyeight/personal"
  "fmt"
)

func RegisterRoutes() {

  // DEFAULT ROUTE
  web.Get("/", func(ctx * web.Context) {
		ctx.Write(toJson( "let thy object decend as if it were calescent" ));
  })

  web.Get("/user/([A-Za-z0-9]+)", func(ctx * web.Context, val string) {
    p:= personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    name := fmt.Sprintf( "@%s", val )
    person, error := p.FindByName( name )
    if error != nil {
      apiError( ctx, error.Error() )
      return
    }

    ctx.Write(toJson( person ))
  })
}
