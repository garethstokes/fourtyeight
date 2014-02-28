package main

import (
  "fmt"
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/garethstokes/fourtyeight/personal"
	"github.com/garethstokes/fourtyeight/cache"
	"github.com/garethstokes/fourtyeight/passwords"
	"github.com/garethstokes/fourtyeight/mail"
)

func WarmAuthCache(){
    // check if user is logged in
    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    p.FillCacheWithLoginTokens()
}

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

    person, error := p.FindByName( val )
    if error != nil {
      apiError( ctx, error.Error() )
      return
    }

    ok( ctx, person )
  })

  web.Get("/user/search/([A-Za-z0-9]+)", func(ctx * web.Context, val string) {
    p:= personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    persons := p.FindByString( val )
    
    ok( ctx, persons )
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

    name := params.Name 
    user, error := p.Validate( name, params.Password )
    if error != nil {
      apiError( ctx, error.Error() )
      return
    }

    hash := passwords.Compute( name + params.Password )

    cache.Set("users", hash.Hash, user )

    user.LoginToken = hash.Hash
    
    error = p.Update(user)
    
    if error != nil {
      apiError( ctx, error.Error() )
      return
    }

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

    user := cache.Get("users", token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    ok( ctx, user )
  })

  web.Get("/users", func(ctx * web.Context) {
    ctx.SetHeader("Context-Type", "application/json", true)

    // NO AUTH
    // user := cache.Get("users", token )
    // if user == nil {
    //   apiError( ctx, "Invalid token" )
    //   return
    // }

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    users := p.FindAll()
    ok( ctx, users )
  })

  web.Get("/user/send/test/email", func(ctx * web.Context) {
    mail.SendWelcomeEmail("gareth@betechnolgy.com.au")
    ok( ctx, "its all good man" )
  });

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
    u.Username    = params.Username
    u.Email       = params.Email
    u.AvatarUrl   = params.AvatarUrl

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
    cache.Set("users", hash.Hash, user)

    go mail.SendWelcomeEmail(u.Email)

    ok( ctx, map[string] interface{} {
      "token": hash.Hash,
      "user": user,
    })
  })

  web.Post("/user/(.+)/follow/(.+)", func(ctx * web.Context, token string, toFollow string) {
    ctx.SetHeader("Context-Type", "application/json", true)

    user := cache.Get("users", token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    personToFollow, _ := p.FindByName( toFollow )
    if personToFollow == nil {
      apiError( ctx, "Invalid follow username." )
      return
    }

    person := user.(* personal.Person)

    //Add to target Followers list
    fmt.Printf("%@ wants to follow %@\n", person.Username, toFollow)

    found := false

    for _, v := range personToFollow.Followers {
      if (v == person.Username) { 
        fmt.Printf("%@ is already a follower of %@\n", person.Username, toFollow)
        found = true
      }
    }
    
    if(!found){
      personToFollow.Followers = append(personToFollow.Followers, person.Username)
      p.Update(personToFollow)
    }

    //reset flag
    found = false

    //Add to target Following list
    for _, v := range person.Following {
      if (v == toFollow) {
        //already following
        fmt.Printf("%@ is already following %@\n", person.Username, toFollow)
        found = true
      }
    }
    
    if(!found){
      person.Following = append(person.Following, toFollow)  
      p.Update(person)
    }
    
    ok( ctx, map[string] interface{} {
      "user": person,
    })

  })

  web.Post("/user/(.+)/unfollow/(.+)", func(ctx * web.Context, token string, toFollow string) {
    ctx.SetHeader("Context-Type", "application/json", true)

    user := cache.Get("users", token )
    if user == nil {
      apiError( ctx, "Invalid token" )
      return
    }

    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()

    personToFollow, _ := p.FindByName( toFollow )
    if personToFollow == nil {
      apiError( ctx, "Invalid follow username." )
      return
    }

    person := user.(* personal.Person)

    //Remove from target Followers list
    fmt.Printf("%@ wants to unfollow %@\n", person.Username, toFollow)

    //target
    followers := make([]string, 0)
    for _, v := range personToFollow.Followers {
      if (v != person.Username) { 
          followers = append(followers, v)
      }
    }
    personToFollow.Followers = followers
    p.Update(personToFollow)
    
    //source
    following := make([]string, 0)
    for _, v := range person.Following {
      if (v != toFollow) {
        following = append(following, v)
      }
    }
    person.Following = following  
    p.Update(person)
    
    ok( ctx, map[string] interface{} {
      "user": person,
    })

  })
}