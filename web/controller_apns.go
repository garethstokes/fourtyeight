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

func registerDevice(user * personal.Person, deviceType int, deviceToken string) {
  // check if token is already in use
  for p, v := range user.NotificationTokens {
      if (v.Token == deviceToken) {
        return
      }
  }

  token := personal.NotificationToken{
    deviceType,
    deviceToken,
  }

  user.NotificationTokens = append(user.NotificationTokens, token)
}

func ApnsController() {


  // register a client for push notifications
  type ApnsRegisterParams struct {
    Token string
    DeviceToken string
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

    if params.DeviceToken == ""{
      apiError( ctx, "INVALID_NOTIFY_TOKEN" )
      fmt.Println("INVALID_NOTIFY_TOKEN")
      return
    }

    if params.Token == "" {
      apiError( ctx, "INVALID_TOKEN" )
      fmt.Println("INVALID_USER_TOKEN")
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

    registerDevice(user, personal.IOS, params.DeviceToken)
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

    registerDevice(user, personal.IOS, params.DeviceToken)
    ok( ctx, params.DeviceToken )
  })
}
