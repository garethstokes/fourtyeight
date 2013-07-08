package main

import (
  "fmt"
	"encoding/json"
	"github.com/garethstokes/web"
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/cache"
	"github.com/garethstokes/fourtyeight/passwords"
)

func RegisterRoutes() {

  // DEFAULT ROUTE
  web.Get("/", func(ctx * web.Context) {
		ctx.Write(toJson( "let thy object decend as if it were calescent" ));
  })

  // GET User by name
  // Example: /user/garrydanger
  //
  // retreives a user
  //
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

  // POST Validate user creds
  // Example: /user/validate
  //
  // params: { name: "garrydanger", password: "bobafett" }
  //
  // will return a session token that the front end can use
  // to access the rest of the api
  type loginParams struct {
    name string
    password string
  }

  web.Post("/user/login", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true);

    params := new( loginParams )
    err := json.NewDecoder(ctx.Request.Body).Decode(&params)
    if err != nil {
			apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
			return
    }

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    name := fmt.Sprintf( "@%s", params.name )
    user, error := p.Validate( name, params.password )
    if error != nil {
      apiError( ctx, error.Error() )
      return
    }

    hash := passwords.Compute( name + params.password )

    cache.Set( hash.Hash, user )

    ctx.Write(toJson(apiOk( map[string] interface{} {
      "token": hash.Hash,
      "user": user,
    })))
  })
}
