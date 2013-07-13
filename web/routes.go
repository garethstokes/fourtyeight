package main

import (
  "fmt"
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/garethstokes/fourtyeight/personal"
)

func RegisterRoutes() {

  // DEFAULT ROUTE
  web.Get("/", func(ctx * web.Context) {
		ctx.Write(toJson( "let thy object decend as if it were calescent" ));
  })

  // Take a email down for the waiting list
  type waitlistParams struct {
    Email string
  }
  web.Post("/waitinglist", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true);

    params := new( waitlistParams )
    err := json.NewDecoder(ctx.Request.Body).Decode(&params)
    if err != nil {
			apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
			return
    }

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    err = p.SaveToWaitingList( params.Email )
    if err != nil {
      apiError( ctx, "a problem happened saving to disk." )
      return
    }

    ok( ctx, "" )
  })

}
