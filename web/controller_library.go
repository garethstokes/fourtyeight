package main

import (
  "fmt"
  "time"
	"encoding/json"
	"github.com/garethstokes/web"
	"github.com/garethstokes/fourtyeight/library"
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/cache"
)

func LibraryController() {

  // GET Documents from the library for a user
  // Example: /library/24332d32e2134231432r
  //
  // will find all relevent documents for a user, ordered
  // by the create_created
  //
  web.Get("/library/(.+)", func(ctx * web.Context, token string) {
    ctx.SetHeader("Content-Type", "application/json", true);

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    user := cache.Get( token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    posts := l.FindAllFor( token )

    ctx.Write(toJson(apiOk( posts )))
  })

  web.Post("/library/(.+)/document", func(ctx * web.Context, token string) {
    ctx.SetHeader("Content-Type", "application/json", true)

    user := cache.Get( token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    post := new( library.Post )
    err := json.NewDecoder(ctx.Request.Body).Decode(&post)
    if err != nil {
			apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
			return
    }

    post.OwnerId = user.(* personal.Person).Username
    post.DateCreated = time.Now().UTC()

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    document := l.CreateFrom( post )

    ctx.Write(toJson(apiOk( document )))
  })
}
