package main

import (
  "fmt"
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/cache"
	"github.com/garethstokes/fourtyeight/passwords"
)

func PersonalController() {

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

    ok( ctx, person )
  })

  // POST Validate user creds
  // Example: /user/validate
  //
  // params: { name: "garrydanger", password: "bobafett" }
  //
  // will return a session token that the front end can use
  // to access the rest of the api
  type loginParams struct {
    Name string
    Password string
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

    name := fmt.Sprintf( "@%s", params.Name )
    user, error := p.Validate( name, params.Password )
    if error != nil {
      apiError( ctx, error.Error() )
      return
    }

    hash := passwords.Compute( name + params.Password )

    cache.Set( hash.Hash, user )

    ok( ctx, map[string] interface{} {
      "token": hash.Hash,
      "user": user,
    })
  })

  // GET the logged in user
  // /me/{token}
  //
  // returns the logged in user for the given token
  web.Get("/me/(.+)", func(ctx * web.Context, token string) {
    ctx.SetHeader("Content-Type", "application/json", true)

    user := cache.Get( token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    ok( ctx, user )
  })

  web.Get("/user/(.+)/following", func(ctx * web.Context, token string) {
    ctx.SetHeader("Context-Type", "application/json", true)

    user := cache.Get( token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    followers, _ := p.Following( user.(* personal.Person) )
    ok( ctx, followers )
  })

  web.Get("/user/(.+)/followers", func(ctx * web.Context, token string) {
    ctx.SetHeader("Context-Type", "application/json", true)

    user := cache.Get( token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    following, _ := p.FollowersFor( user.(* personal.Person) )
    ok( ctx, following )
  })

  web.Post("/user/(.+)/follow/(.+)", func(ctx * web.Context, token string, toFollow string) {
    ctx.SetHeader("Context-Type", "application/json", true)

    user := cache.Get( token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    personToFollow, _ := p.FindByName(fmt.Sprintf( "@%s", toFollow ))
    if personToFollow == nil {
      apiError( ctx, "Invalid follow username." )
      return
    }

    following, _ := p.AddFollowerTo( personToFollow, user.(* personal.Person) )
    ok( ctx, following )
  })

  /*
      CREATE USER: 
        note that ctx.Params does not work on json input

  */
  type userCreateParams struct {
    personal.Person
    Password string
  }

  web.Post("/user", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true);

    params := new(userCreateParams)
    err := json.NewDecoder(ctx.Request.Body).Decode(&params)
    if err != nil {
			apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
			return
    }

    fmt.Printf("%@\n", params)

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    u := new(personal.Person)
    u.Username = fmt.Sprintf( "@%s", params.Username )
    u.Email = params.Email
    u.AvatarUrl = params.AvatarUrl
    fmt.Printf("%@\n", u)

    validations := u.Validate()
    if len(validations) > 0 {
      fmt.Println(validations)
      apiError(ctx, validations)
      return
    }

    user, error := p.Create(u, params.Password)
    if error != nil {
      apiError( ctx, error.Error() )
      return
    }

    hash := passwords.Compute( user.Username + params.Password )
    cache.Set(hash.Hash, user)

    ok( ctx, map[string] interface{} {
      "token": hash.Hash,
      "user": user,
    })
  })

}
