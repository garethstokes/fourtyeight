package main

import (
  "fmt"
	"encoding/json"
	"github.com/hoisie/web"
  "github.com/garethstokes/fourtyeight/cache"
	"github.com/garethstokes/fourtyeight/personal"
)

func WarmApnCache(){
    // check if user is logged in
    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()
    
    p.FillCacheWithNotificationTokens()

}

func ApnsController() {

  // register a client for push notifications
  type ApnsRegisterParams struct {
    Token string
    DeviceToken string
    DeviceType string
  }

  web.Post("/apns/register", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true)

    // get params
    params := new(ApnsRegisterParams)
    err := json.NewDecoder(ctx.Request.Body).Decode(&params)
    if err != nil {
      apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
      return
    }

    // check if user is logged in
    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()
    
    user, error := p.GetLoggedInUser(params.Token)

    if error != nil {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN")
      return
    }

    //woo hoo it worked, we found them

    //put it in the cache i guess
    cache.Set("apns", user.Username, params.DeviceToken)
    
    //put it in the database is better
    error = p.RegisteriOSDevice(user.Username, params.DeviceToken)

    if error != nil {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN - failed to store the device token for apns")

      return
    }

    ok( ctx, params.DeviceToken )
  })


  web.Post("/apns/register/android", func(ctx * web.Context) {
    ctx.SetHeader("Content-Type", "application/json", true)

    // get params
    params := new(ApnsRegisterParams)
    err := json.NewDecoder(ctx.Request.Body).Decode(&params)
    if err != nil {
      apiError(ctx, "incorrect parameters found")
      fmt.Printf( "ERROR: %s\n", err.Error() )
      return
    }

    // check if user is logged in
    p := personal.Store()
    p.OpenSession()
    defer p.CloseSession()
    
    user, error := p.GetLoggedInUser(params.Token)

    if error != nil {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN")
      return
    }

    //woo hoo it worked, we found them

    //put it in the cache i guess
    cache.Set("apns_android", user.Username, params.DeviceToken)
    
    //put it in the database is better
    error = p.RegisterAndroidDevice(user.Username, params.DeviceToken)

    if error != nil {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN - failed to store the device token for apns")

      return
    }

    ok( ctx, params.DeviceToken )
  })
}
