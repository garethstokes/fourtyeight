package main

import (
  "fmt"
	"encoding/json"
	"github.com/hoisie/web"
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

    posts := l.FindDocumentsFor( token )

    ok( ctx, posts )
  })

  // POST Create a document attached to a user
  // Example: localhost:8080/library/e51n4EZvN8KL7uoQUtmbWw==/document 
  //
  // will attach a expiry delta that needs to be specified
  //
  type PostWithExpiry struct {
    library.Post
    Expiry int64
  }
  web.Post("/library/(.+)/document", func(ctx * web.Context, token string) {
    ctx.SetHeader("Content-Type", "application/json", true)

    user := cache.Get( token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    post := new(PostWithExpiry)
    err := json.NewDecoder(ctx.Request.Body).Decode(&post)
    if err != nil {
			apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
			return
    }

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    p := new(library.Post)
    p.OwnerId = user.(* personal.Person).Username
    p.Image = post.Image
    p.Text = post.Text

    document := l.CreateFrom( p, post.Expiry )

    ok( ctx, document )
  })

  web.Get("/document/(.+)", func(ctx * web.Context, documentId string) {
    ctx.SetHeader("Content-Type", "application/json", true)

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    document := l.FindOne( documentId )
    if document == nil {
      apiError( ctx, "Incorrect document id" )
      return
    }

    ok( ctx, document )
  })

  web.Post("/library/(.+)/document/(.+)/post", func(ctx * web.Context, token string, documentId string) {
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

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    document := l.AddPost( post, documentId )

    ok( ctx, document )
  })

  web.Post("/library/(.+)/delete/(.+)", func(ctx * web.Context, token string, documentId string) {
    ctx.SetHeader("Content-Type", "application/json", true)

    user := cache.Get(token)
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    l := library.Store()
    l.OpenSession()
    defer l.CloseSession()

    err := l.DeleteOne(documentId)
    if err != nil {
      apiError(ctx, "There was a problem contacting library service.")
      fmt.Printf("ERROR: %s\n", err.Error())
      return
    }

    ok( ctx, true )
  })
}
